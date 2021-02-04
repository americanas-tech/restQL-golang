package web

import (
	"github.com/b2wdigital/restQL-golang/v4/internal/platform/persistence"
	"github.com/b2wdigital/restQL-golang/v4/pkg/restql"
	"github.com/valyala/fasthttp"
)

type administrator struct {
	mr persistence.MappingsReader
	qr persistence.QueryReader
}

func newAdmin(mr persistence.MappingsReader, qr persistence.QueryReader) *administrator {
	return &administrator{mr: mr, qr: qr}
}

func (adm *administrator) ListAllTenants(ctx *fasthttp.RequestCtx) error {
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

func (adm administrator) ListAllNamespaces(ctx *fasthttp.RequestCtx) error {
	namespaces, err := adm.qr.ListNamespaces(ctx)
	if err != nil {
		return RespondError(ctx, err)
	}

	data := map[string]interface{}{"namespaces": namespaces}
	return Respond(ctx, data, fasthttp.StatusOK, nil)
}

type query struct {
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
		return err
	}

	queries := make(map[string][]query)
	for queryName, savedQueries := range queriesForNamespace {
		qs := make([]query, len(savedQueries))
		for i, savedQuery := range savedQueries {
			qs[i] = query{
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
