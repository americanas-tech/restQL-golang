package web

import (
	"github.com/b2wdigital/restQL-golang/v4/internal/platform/persistence"
	"github.com/b2wdigital/restQL-golang/v4/pkg/restql"
	"github.com/valyala/fasthttp"
	"strconv"
)

type administrator struct {
	mr persistence.MappingsReader
	qr persistence.QueryReader
}

func newAdmin(mr persistence.MappingsReader, qr persistence.QueryReader) *administrator {
	return &administrator{mr: mr, qr: qr}
}

func (adm *administrator) AllTenants(ctx *fasthttp.RequestCtx) error {
	tenants, err := adm.mr.ListTenants(ctx)
	if err != nil {
		return RespondError(ctx, err)
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

	mappings, err := adm.mr.FromTenant(ctx, tenantName)
	if err != nil {
		return RespondError(ctx, err)
	}

	urls := make(map[string]string)
	for resourceName, mapping := range mappings {
		urls[resourceName] = mapping.URL()
	}

	data := map[string]interface{}{
		"tenant":   tenantName,
		"mappings": urls,
	}
	return Respond(ctx, data, fasthttp.StatusOK, nil)
}

func (adm administrator) AllNamespaces(ctx *fasthttp.RequestCtx) error {
	namespaces, err := adm.qr.ListNamespaces(ctx)
	if err != nil {
		return RespondError(ctx, err)
	}

	data := map[string]interface{}{"namespaces": namespaces}
	return Respond(ctx, data, fasthttp.StatusOK, nil)
}

type queryRevision struct {
	Name     string `json:"name"`
	Text     string `json:"text"`
	Revision int    `json:"revision"`
}

func (adm administrator) NamespaceQueries(ctx *fasthttp.RequestCtx) error {
	log := restql.GetLogger(ctx)

	namespace, err := pathParamString(ctx, "namespace")
	if err != nil {
		log.Error("failed to load namespace path param", err)
		return err
	}

	queriesForNamespace, err := adm.qr.ListQueriesForNamespace(ctx, namespace)
	if err != nil {
		return RespondError(ctx, err)
	}

	queries := make(map[string][]queryRevision)
	for queryName, savedQueries := range queriesForNamespace {
		qs := make([]queryRevision, len(savedQueries))
		for i, savedQuery := range savedQueries {
			qs[i] = queryRevision{
				Name:     savedQuery.Name,
				Text:     savedQuery.Text,
				Revision: savedQuery.Revision,
			}
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
		log.Error("failed to load queryRevision name path param", err)
		return err
	}

	rs, err := adm.qr.ListQueryRevisions(ctx, namespace, queryName)
	if err != nil {
		return RespondError(ctx, err)
	}

	queryRevisions := make([]queryRevision, len(rs))
	for i, r := range rs {
		queryRevisions[i] = queryRevision{
			Name:     r.Name,
			Text:     r.Text,
			Revision: r.Revision,
		}
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
		log.Error("failed to load queryRevision name path param", err)
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
		return RespondError(ctx, err)
	}

	data := map[string]interface{}{
		"namespace": namespace,
		"name":      savedQuery.Name,
		"revision": map[string]string{
			"text": savedQuery.Text,
		},
	}

	return Respond(ctx, data, fasthttp.StatusOK, nil)
}
