package web

import (
	"bytes"
	"encoding/json"
	"github.com/b2wdigital/restQL-golang/v5/internal/platform/persistence"
	"github.com/b2wdigital/restQL-golang/v5/internal/platform/web/middleware"
	"github.com/b2wdigital/restQL-golang/v5/pkg/restql"
	"github.com/valyala/fasthttp"
	"strconv"
)

type queryRevision struct {
	Name     string `json:"name,omitempty"`
	Text     string `json:"text,omitempty"`
	Revision int    `json:"revision,omitempty"`
	Source   string `json:"source,omitempty"`
}

type query struct {
	Namespace string          `json:"namespace"`
	Name      string          `json:"name"`
	Revisions []queryRevision `json:"revisions"`
	Source    string          `json:"source,omitempty"`
}

type mapping struct {
	URL    string `json:"url"`
	Source string `json:"source"`
}

type administrator struct {
	log               restql.Logger
	mr                persistence.MappingsReader
	mw                persistence.MappingsWriter
	qr                persistence.QueryReader
	queryWriter       persistence.QueryWriter
	authorizationCode []byte
}

func newAdmin(log restql.Logger, mr persistence.MappingsReader, mw persistence.MappingsWriter, qr persistence.QueryReader, qw persistence.QueryWriter, authorizationCode string) *administrator {
	return &administrator{log: log, mr: mr, mw: mw, qr: qr, queryWriter: qw, authorizationCode: []byte(authorizationCode)}
}

func (adm *administrator) AllTenants(ctx *fasthttp.RequestCtx) error {
	tenants, err := adm.mr.ListTenants(ctx)
	if err != nil {
		return RespondError(ctx, err, errToStatusCode)
	}

	data := map[string]interface{}{"tenants": tenants}
	return Respond(ctx, data, fasthttp.StatusOK, nil)
}

func (adm *administrator) TenantMappings(reqCtx *fasthttp.RequestCtx) error {
	ctx := middleware.GetNativeContext(reqCtx)
	ctx = restql.WithLogger(ctx, adm.log)

	tenantName, err := pathParamString(reqCtx, "tenantName")
	if err != nil {
		adm.log.Error("failed to load tenant name path param", err)
		return err
	}

	sourceFilter := restql.Source(reqCtx.QueryArgs().Peek("source"))

	mappings, err := adm.mr.FromTenant(ctx, tenantName)
	if err != nil {
		return RespondError(reqCtx, err, errToStatusCode)
	}

	mappings = filterMappingsBySource(mappings, sourceFilter)

	ms := make(map[string]mapping)
	for resourceName, m := range mappings {
		ms[resourceName] = mapping{
			URL:    m.URL(),
			Source: string(m.Source),
		}
	}

	data := map[string]interface{}{
		"tenant":   tenantName,
		"mappings": ms,
	}
	return Respond(reqCtx, data, fasthttp.StatusOK, nil)
}

func (adm administrator) AllNamespaces(reqCtx *fasthttp.RequestCtx) error {
	ctx := middleware.GetNativeContext(reqCtx)
	ctx = restql.WithLogger(ctx, adm.log)

	namespaces, err := adm.qr.ListNamespaces(ctx)
	if err != nil {
		return RespondError(reqCtx, err, errToStatusCode)
	}

	data := map[string]interface{}{"namespaces": namespaces}
	return Respond(reqCtx, data, fasthttp.StatusOK, nil)
}

func (adm administrator) NamespaceQueries(reqCtx *fasthttp.RequestCtx) error {
	ctx := middleware.GetNativeContext(reqCtx)
	ctx = restql.WithLogger(ctx, adm.log)

	namespace, err := pathParamString(reqCtx, "namespace")
	if err != nil {
		adm.log.Error("failed to load namespace path param", err)
		return err
	}

	sourceFilter := restql.Source(reqCtx.QueryArgs().Peek("source"))

	queriesForNamespace, err := adm.qr.ListQueriesForNamespace(ctx, namespace)
	if err != nil {
		return RespondError(reqCtx, err, errToStatusCode)
	}

	queries := make([]query, len(queriesForNamespace))
	i := 0
	for queryName, savedQueries := range queriesForNamespace {
		savedQueries = filterQueriesBySource(savedQueries, sourceFilter)

		rs := make([]queryRevision, len(savedQueries))
		for j, savedQuery := range savedQueries {
			rs[j] = queryRevision{
				Text:     savedQuery.Text,
				Revision: savedQuery.Revision,
				Source:   string(savedQuery.Source),
			}
		}

		queries[i] = query{Name: queryName, Namespace: namespace, Revisions: rs}
		i++
	}

	data := map[string]interface{}{"namespace": namespace, "queries": queries}
	return Respond(reqCtx, data, fasthttp.StatusOK, nil)
}

func (adm *administrator) QueryRevisions(reqCtx *fasthttp.RequestCtx) error {
	ctx := middleware.GetNativeContext(reqCtx)
	ctx = restql.WithLogger(ctx, adm.log)

	namespace, err := pathParamString(reqCtx, "namespace")
	if err != nil {
		adm.log.Error("failed to load namespace path param", err)
		return err
	}

	queryName, err := pathParamString(reqCtx, "queryId")
	if err != nil {
		adm.log.Error("failed to load query name path param", err)
		return err
	}

	sourceFilter := restql.Source(reqCtx.QueryArgs().Peek("source"))

	rs, err := adm.qr.ListQueryRevisions(ctx, namespace, queryName)
	if err != nil {
		return RespondError(reqCtx, err, errToStatusCode)
	}

	rs = filterQueriesBySource(rs, sourceFilter)

	queryRevisions := make([]queryRevision, len(rs))
	for i, r := range rs {
		queryRevisions[i] = toQueryRevision(r)
	}

	data := query{
		Namespace: namespace,
		Name:      queryName,
		Revisions: queryRevisions,
	}
	return Respond(reqCtx, data, fasthttp.StatusOK, nil)
}

func (adm *administrator) Query(reqCtx *fasthttp.RequestCtx) error {
	ctx := middleware.GetNativeContext(reqCtx)
	ctx = restql.WithLogger(ctx, adm.log)

	namespace, err := pathParamString(reqCtx, "namespace")
	if err != nil {
		adm.log.Error("failed to load namespace path param", err)
		return err
	}

	queryName, err := pathParamString(reqCtx, "queryId")
	if err != nil {
		adm.log.Error("failed to load query name path param", err)
		return err
	}

	revisionStr, err := pathParamString(reqCtx, "revision")
	if err != nil {
		adm.log.Error("failed to load revision path param", err)
		return err
	}

	revision, err := strconv.Atoi(revisionStr)
	if err != nil {
		adm.log.Error("failed to parse revision path param", err)
		return err
	}

	savedQuery, err := adm.qr.Get(ctx, namespace, queryName, revision)
	if err != nil {
		return RespondError(reqCtx, err, errToStatusCode)
	}

	data := map[string]interface{}{
		"namespace": namespace,
		"name":      savedQuery.Name,
		"source":    savedQuery.Source,
		"revision": map[string]string{
			"text": savedQuery.Text,
		},
	}

	return Respond(reqCtx, data, fasthttp.StatusOK, nil)
}

type mapResourceBody struct {
	Url string `json:"url"`
}

func (adm *administrator) MapResource(reqCtx *fasthttp.RequestCtx) error {
	ctx := middleware.GetNativeContext(reqCtx)
	ctx = restql.WithLogger(ctx, adm.log)

	if !isAuthorized(reqCtx, adm.authorizationCode) {
		reqCtx.Response.SetStatusCode(fasthttp.StatusUnauthorized)
		return nil
	}

	tenantName, err := pathParamString(reqCtx, "tenantName")
	if err != nil {
		adm.log.Error("failed to load tenant name path param", err)
		return err
	}

	resourceName, err := pathParamString(reqCtx, "resource")
	if err != nil {
		adm.log.Error("failed to load resource name path param", err)
		return err
	}

	var mrb mapResourceBody

	bytesBody := reqCtx.PostBody()
	err = json.Unmarshal(bytesBody, &mrb)
	if err != nil {
		return err
	}

	err = adm.mw.Write(ctx, tenantName, resourceName, mrb.Url)
	if err != nil {
		return RespondError(reqCtx, err, errToStatusCode)
	}

	return Respond(reqCtx, nil, fasthttp.StatusCreated, nil)
}

func isAuthorized(ctx *fasthttp.RequestCtx, authorizationCode []byte) bool {
	bearerCode := getBearerToken(ctx)
	if len(bearerCode) == 0 {
		return false
	}

	bearerCode = bytes.TrimPrefix(bearerCode, []byte("Bearer"))
	bearerCode = bytes.TrimPrefix(bearerCode, []byte("bearer"))
	bearerCode = bytes.TrimSpace(bearerCode)

	return bytes.Equal(bearerCode, authorizationCode)
}

func getBearerToken(ctx *fasthttp.RequestCtx) []byte {
	bearerCode := ctx.Request.Header.Peek("Authorization")
	if len(bearerCode) > 0 {
		return bearerCode
	}

	bearerCode = ctx.Request.Header.Peek("authorization")

	return bearerCode
}

type createRevisionBody struct {
	Text string `json:"text"`
}

func (adm *administrator) CreateQueryRevision(reqCtx *fasthttp.RequestCtx) error {
	ctx := middleware.GetNativeContext(reqCtx)
	ctx = restql.WithLogger(ctx, adm.log)

	namespace, err := pathParamString(reqCtx, "namespace")
	if err != nil {
		adm.log.Error("failed to load namespace path param", err)
		return err
	}

	queryName, err := pathParamString(reqCtx, "queryId")
	if err != nil {
		adm.log.Error("failed to load query name path param", err)
		return err
	}

	var crb createRevisionBody

	bytesBody := reqCtx.PostBody()
	err = json.Unmarshal(bytesBody, &crb)
	if err != nil {
		return err
	}

	err = adm.queryWriter.Write(ctx, namespace, queryName, crb.Text)
	if err != nil {
		return RespondError(reqCtx, err, errToStatusCode)
	}

	return Respond(reqCtx, nil, fasthttp.StatusCreated, nil)
}

func filterQueriesBySource(queryRevisions []restql.SavedQuery, source restql.Source) []restql.SavedQuery {
	if source == "" {
		return queryRevisions
	}

	var result []restql.SavedQuery

	for _, qr := range queryRevisions {
		if qr.Source == source {
			result = append(result, qr)
		}
	}

	return result
}

func filterMappingsBySource(mappings map[string]restql.Mapping, source restql.Source) map[string]restql.Mapping {
	if source == "" {
		return mappings
	}

	result := make(map[string]restql.Mapping)

	for _, m := range mappings {
		if m.Source == source {
			result[m.ResourceName()] = m
		}
	}

	return result
}

func toQueryRevision(sq restql.SavedQuery) queryRevision {
	return queryRevision{
		Name:     sq.Name,
		Text:     sq.Text,
		Revision: sq.Revision,
		Source:   string(sq.Source),
	}
}
