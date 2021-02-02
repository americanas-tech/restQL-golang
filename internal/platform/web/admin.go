package web

import (
	"github.com/b2wdigital/restQL-golang/v4/internal/platform/persistence"
	"github.com/valyala/fasthttp"
)

type administrator struct {
	admRepo *persistence.AdminRepository
}

func newAdmin(admRepo *persistence.AdminRepository) *administrator {
	return &administrator{admRepo: admRepo}
}

func (adm *administrator) ListAllTenants(ctx *fasthttp.RequestCtx) error {
	tenants, err := adm.admRepo.ListAllTenants(ctx)
	if err != nil {
		return RespondError(ctx, err)
	}

	data := map[string]interface{}{"tenants": tenants}
	return Respond(ctx, data, fasthttp.StatusOK, nil)
}
