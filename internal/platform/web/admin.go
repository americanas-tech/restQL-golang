package web

import (
	"github.com/b2wdigital/restQL-golang/v4/internal/platform/persistence"
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
