package web

import (
	"github.com/b2wdigital/restQL-golang/v4/internal/platform/persistence"
	"github.com/b2wdigital/restQL-golang/v4/pkg/restql"
	"github.com/valyala/fasthttp"
)

type administrator struct {
	mr persistence.MappingsReader
}

func newAdmin(mr persistence.MappingsReader) *administrator {
	return &administrator{mr: mr}
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
