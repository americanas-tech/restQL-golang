package web

import (
	"encoding/json"
	"github.com/b2wdigital/restQL-golang/v4/internal/platform/persistence"
	"github.com/b2wdigital/restQL-golang/v4/pkg/restql"
	"github.com/valyala/fasthttp"
	"strconv"
)

type queryRevision struct {
	Name     string `json:"name"`
	Text     string `json:"text"`
	Revision int    `json:"revision"`
	Source   string `json:"source"`
}

type mapping struct {
	URL    string `json:"url"`
	Source string `json:"source"`
}

type administrator struct {
	mr          persistence.MappingsReader
	mw          persistence.MappingsWriter
	qr          persistence.QueryReader
	queryWriter persistence.QueryWriter
}

func newAdmin(mr persistence.MappingsReader, mw persistence.MappingsWriter, qr persistence.QueryReader, qw persistence.QueryWriter) *administrator {
	return &administrator{mr: mr, mw: mw, qr: qr, queryWriter: qw}
}

func (adm *administrator) AllTenants(ctx *fasthttp.RequestCtx) error {
	tenants, err := adm.mr.ListTenants(ctx)
	if err != nil {
		return RespondError(ctx, err, errToStatusCode)
	}

	data := map[string]interface{}{"tenants": tenants}
	return Respond(ctx, data, fasthttp.StatusOK, nil)
}

func (adm *administrator) TenantMappings(ctx *fasthttp.RequestCtx) error {
	log := restql.GetLogger(ctx)

	tenantName, err := pathParamString(ctx, "tenantName")
	if err != nil {
		log.Error("failed to load tenant name path param", err)
		return err
	}

	sourceFilter := restql.Source(ctx.QueryArgs().Peek("source"))

	mappings, err := adm.mr.FromTenant(ctx, tenantName)
	if err != nil {
		return RespondError(ctx, err, errToStatusCode)
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
	return Respond(ctx, data, fasthttp.StatusOK, nil)
}

func (adm administrator) AllNamespaces(ctx *fasthttp.RequestCtx) error {
	namespaces, err := adm.qr.ListNamespaces(ctx)
	if err != nil {
		return RespondError(ctx, err, errToStatusCode)
	}

	data := map[string]interface{}{"namespaces": namespaces}
	return Respond(ctx, data, fasthttp.StatusOK, nil)
}

func (adm administrator) NamespaceQueries(ctx *fasthttp.RequestCtx) error {
	log := restql.GetLogger(ctx)

	namespace, err := pathParamString(ctx, "namespace")
	if err != nil {
		log.Error("failed to load namespace path param", err)
		return err
	}

	sourceFilter := restql.Source(ctx.QueryArgs().Peek("source"))

	queriesForNamespace, err := adm.qr.ListQueriesForNamespace(ctx, namespace)
	if err != nil {
		return RespondError(ctx, err, errToStatusCode)
	}

	queries := make(map[string][]queryRevision)
	for queryName, savedQueries := range queriesForNamespace {
		savedQueries = filterQueriesBySource(savedQueries, sourceFilter)
		if len(savedQueries) == 0 {
			continue
		}

		qs := make([]queryRevision, len(savedQueries))
		for i, savedQuery := range savedQueries {
			qs[i] = toQueryRevision(savedQuery)
		}

		queries[queryName] = qs
	}

	data := map[string]interface{}{"namespace": namespace, "queries": queries}
	return Respond(ctx, data, fasthttp.StatusOK, nil)
}

func (adm *administrator) QueryRevisions(ctx *fasthttp.RequestCtx) error {
	log := restql.GetLogger(ctx)

	namespace, err := pathParamString(ctx, "namespace")
	if err != nil {
		log.Error("failed to load namespace path param", err)
		return err
	}

	queryName, err := pathParamString(ctx, "queryId")
	if err != nil {
		log.Error("failed to load query name path param", err)
		return err
	}

	sourceFilter := restql.Source(ctx.QueryArgs().Peek("source"))

	rs, err := adm.qr.ListQueryRevisions(ctx, namespace, queryName)
	if err != nil {
		return RespondError(ctx, err, errToStatusCode)
	}

	rs = filterQueriesBySource(rs, sourceFilter)

	queryRevisions := make([]queryRevision, len(rs))
	for i, r := range rs {
		queryRevisions[i] = toQueryRevision(r)
	}

	data := map[string]interface{}{"namespace": namespace, "query": queryName, "revisions": queryRevisions}
	return Respond(ctx, data, fasthttp.StatusOK, nil)
}

func (adm *administrator) Query(ctx *fasthttp.RequestCtx) error {
	log := restql.GetLogger(ctx)

	namespace, err := pathParamString(ctx, "namespace")
	if err != nil {
		log.Error("failed to load namespace path param", err)
		return err
	}

	queryName, err := pathParamString(ctx, "queryId")
	if err != nil {
		log.Error("failed to load query name path param", err)
		return err
	}

	revisionStr, err := pathParamString(ctx, "revision")
	if err != nil {
		log.Error("failed to load revision path param", err)
		return err
	}

	revision, err := strconv.Atoi(revisionStr)
	if err != nil {
		log.Error("failed to parse revision path param", err)
		return err
	}

	savedQuery, err := adm.qr.Get(ctx, namespace, queryName, revision)
	if err != nil {
		return RespondError(ctx, err, errToStatusCode)
	}

	data := map[string]interface{}{
		"namespace": namespace,
		"name":      savedQuery.Name,
		"source":    savedQuery.Source,
		"revision": map[string]string{
			"text": savedQuery.Text,
		},
	}

	return Respond(ctx, data, fasthttp.StatusOK, nil)
}

type mapResourceBody struct {
	Url string `json:"url"`
}

func (adm *administrator) MapResource(ctx *fasthttp.RequestCtx) error {
	log := restql.GetLogger(ctx)

	tenantName, err := pathParamString(ctx, "tenantName")
	if err != nil {
		log.Error("failed to load tenant name path param", err)
		return err
	}

	resourceName, err := pathParamString(ctx, "resource")
	if err != nil {
		log.Error("failed to load resource name path param", err)
		return err
	}

	var mrb mapResourceBody

	bytesBody := ctx.PostBody()
	err = json.Unmarshal(bytesBody, &mrb)
	if err != nil {
		return err
	}

	err = adm.mw.Write(ctx, tenantName, resourceName, mrb.Url)
	if err != nil {
		return RespondError(ctx, err, errToStatusCode)
	}

	return Respond(ctx, nil, fasthttp.StatusCreated, nil)
}

type createRevisionBody struct {
	Text string `json:"text"`
}

func (adm *administrator) CreateQueryRevision(ctx *fasthttp.RequestCtx) error {
	log := restql.GetLogger(ctx)

	namespace, err := pathParamString(ctx, "namespace")
	if err != nil {
		log.Error("failed to load namespace path param", err)
		return err
	}

	queryName, err := pathParamString(ctx, "queryId")
	if err != nil {
		log.Error("failed to load query name path param", err)
		return err
	}

	var crb createRevisionBody

	bytesBody := ctx.PostBody()
	err = json.Unmarshal(bytesBody, &crb)
	if err != nil {
		return err
	}

	err = adm.queryWriter.Write(ctx, namespace, queryName, crb.Text)
	if err != nil {
		return RespondError(ctx, err, errToStatusCode)
	}

	return Respond(ctx, nil, fasthttp.StatusCreated, nil)
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
