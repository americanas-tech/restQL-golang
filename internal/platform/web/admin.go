package web

import (
	"bytes"
	"encoding/json"
	"github.com/b2wdigital/restQL-golang/v6/internal/platform/persistence"
	"github.com/b2wdigital/restQL-golang/v6/internal/platform/web/middleware"
	"github.com/b2wdigital/restQL-golang/v6/pkg/restql"
	"github.com/valyala/fasthttp"
	"strconv"
)

type queryRevision struct {
	Text     string `json:"text,omitempty"`
	Revision int    `json:"revision,omitempty"`
	Archived bool   `json:"archived"`
	Source   string `json:"source,omitempty"`
}

type query struct {
	Namespace string          `json:"namespace"`
	Name      string          `json:"name"`
	Revisions []queryRevision `json:"revisions"`
	Archived  bool            `json:"archived"`
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
	archivedFilter, err := strconv.ParseBool(string(reqCtx.QueryArgs().Peek("archived")))
	if err != nil {
		archivedFilter = false
	}

	queriesForNamespace, err := adm.qr.ListQueriesForNamespace(ctx, namespace, archivedFilter)
	if err != nil {
		return RespondError(reqCtx, err, errToStatusCode)
	}

	//queries := make([]query, len(queriesForNamespace))
	queries := []query{}
	for _, savedQuery := range queriesForNamespace {
		filteredRevisions := filterRevisionsBySource(savedQuery, sourceFilter)
		if len(filteredRevisions) == 0 && len(savedQuery.Revisions) != 0 {
			continue
		}

		rs := make([]queryRevision, len(filteredRevisions))
		for j, rev := range filteredRevisions {
			rs[j] = queryRevision{
				Text:     rev.Text,
				Revision: rev.Revision,
				Archived: rev.Archived,
				Source:   string(rev.Source),
			}
		}

		adm.log.Debug("saved query", "value", savedQuery)
		queries = append(queries, query{Name: savedQuery.Name, Namespace: savedQuery.Namespace, Archived: savedQuery.Archived, Revisions: rs})
	}

	adm.log.Debug("namespace queries", "value", queries)

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
	archivedFilter, err := strconv.ParseBool(string(reqCtx.QueryArgs().Peek("archived")))
	if err != nil {
		archivedFilter = false
	}

	savedQuery, err := adm.qr.ListQueryRevisions(ctx, namespace, queryName, archivedFilter)
	if err != nil {
		return RespondError(reqCtx, err, errToStatusCode)
	}

	filteredRevisions := filterRevisionsBySource(savedQuery, sourceFilter)

	queryRevisions := make([]queryRevision, len(filteredRevisions))
	for i, r := range filteredRevisions {
		queryRevisions[i] = toQueryRevision(r)
	}

	data := query{
		Namespace: savedQuery.Namespace,
		Name:      savedQuery.Name,
		Archived:  savedQuery.Archived,
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

	qr := toQueryRevision(savedQuery)

	return Respond(reqCtx, qr, fasthttp.StatusOK, nil)
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

	var crb createRevisionBody
	err := json.Unmarshal(reqCtx.PostBody(), &crb)
	if err != nil {
		return err
	}

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

	err = adm.queryWriter.Write(ctx, namespace, queryName, crb.Text)
	if err != nil {
		return RespondError(reqCtx, err, errToStatusCode)
	}

	return Respond(reqCtx, nil, fasthttp.StatusCreated, nil)
}

type updateArchivingBody struct {
	Archived bool `json:"archived"`
}

func (adm *administrator) UpdateQueryArchiving(reqCtx *fasthttp.RequestCtx) error {
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

	var body updateArchivingBody

	bytesBody := reqCtx.PostBody()
	err = json.Unmarshal(bytesBody, &body)
	if err != nil {
		return err
	}

	err = adm.queryWriter.UpdateQueryArchiving(ctx, namespace, queryName, body.Archived)
	if err != nil {
		return RespondError(reqCtx, err, errToStatusCode)
	}

	return Respond(reqCtx, nil, fasthttp.StatusNoContent, nil)
}

func (adm *administrator) UpdateRevisionArchiving(reqCtx *fasthttp.RequestCtx) error {
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

	var body updateArchivingBody

	bytesBody := reqCtx.PostBody()
	err = json.Unmarshal(bytesBody, &body)
	if err != nil {
		return err
	}

	err = adm.queryWriter.UpdateRevisionArchiving(ctx, namespace, queryName, revision, body.Archived)
	if err != nil {
		return RespondError(reqCtx, err, errToStatusCode)
	}

	return Respond(reqCtx, nil, fasthttp.StatusNoContent, nil)
}

func filterRevisionsBySource(query restql.SavedQuery, source restql.Source) []restql.SavedQueryRevision {
	if source == "" {
		return query.Revisions
	}

	revs := []restql.SavedQueryRevision{}

	for _, r := range query.Revisions {
		if r.Source == source {
			revs = append(revs, r)
		}
	}

	return revs
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

func toQueryRevision(sq restql.SavedQueryRevision) queryRevision {
	return queryRevision{
		Text:     sq.Text,
		Revision: sq.Revision,
		Archived: sq.Archived,
		Source:   string(sq.Source),
	}
}
