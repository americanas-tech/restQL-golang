package persistence

import (
	"context"
	"github.com/b2wdigital/restQL-golang/v4/internal/domain"
	"github.com/b2wdigital/restQL-golang/v4/pkg/restql"
)

type AdminRepository struct {
	log          restql.Logger
	localQueries map[string]map[string][]string
	env          domain.EnvSource
	db           Database
}

func NewAdminRepository(log restql.Logger, db Database) *AdminRepository {
	return &AdminRepository{log: log, db: db}
}

func (adm *AdminRepository) ListAllTenants(ctx context.Context) ([]string, error) {
	tenants, err := adm.db.FindAllTenants(ctx)
	if err != nil {
		return nil, err
	}

	return tenants, nil
}
