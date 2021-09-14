package persistence

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/b2wdigital/restQL-golang/v6/internal/domain"
	"github.com/b2wdigital/restQL-golang/v6/pkg/restql"
)

var envMappingWithTenantRegex = regexp.MustCompile(`^RESTQL_MAPPING_(\w+)_(\w+)$`)

// MappingsReader fetch indexed mappings from database,
// configuration file and environment variable.
type MappingsReader struct {
	log           restql.Logger
	env           map[string]map[string]restql.Mapping
	localByTenant map[string]map[string]restql.Mapping
	db            Database
}

// NewMappingReader constructs a MappingsReader instance.
func NewMappingReader(log restql.Logger, env domain.EnvSource, local map[string]map[string]string, db Database) MappingsReader {
	envWithTenantMappings := getMappingsFromEnv(log, env)
	localMappings := make(map[string]map[string]restql.Mapping)
	for t, m := range local {
		localMappings[t] = parseMappingsFromLocal(log, m)
	}

	return MappingsReader{log: log, env: envWithTenantMappings, localByTenant: localMappings, db: db}
}

// ListTenants fetch all tenants under which mappings are organized
func (mr MappingsReader) ListTenants(ctx context.Context) ([]string, error) {
	tenantSet := make(map[string]struct{})

	for tenant := range mr.localByTenant {
		tenantSet[tenant] = struct{}{}
	}

	dbTenants, err := mr.db.FindAllTenants(ctx)
	if err != nil {
		log := restql.GetLogger(ctx)
		log.Error("fail to read tenants from database", err)
	} else {
		for _, tenant := range dbTenants {
			tenantSet[tenant] = struct{}{}
		}
	}

	for tenant := range mr.env {
		tenantSet[tenant] = struct{}{}
	}

	tenants := make([]string, len(tenantSet))
	i := 0
	for tenant := range tenantSet {
		tenants[i] = tenant
		i++
	}

	return tenants, nil
}

// FromTenant fetch the mappings for the given tenant.
func (mr MappingsReader) FromTenant(ctx context.Context, tenant string) (map[string]restql.Mapping, error) {
	log := restql.GetLogger(ctx)
	log.Debug("fetching mappings")
	errMappingsFound := fmt.Errorf("%w: tenant %s", restql.ErrMappingsNotFound, tenant)

	result := make(map[string]restql.Mapping)

	localTenantMappings, found := mr.localByTenant[tenant]
	if found {
		for k, v := range localTenantMappings {
			result[k] = v
		}
	}

	dbMappings, err := mr.db.FindMappingsForTenant(ctx, tenant)
	switch {
	case err == errNoDatabase || errors.Is(err, restql.ErrMappingsNotFoundInDatabase):
		result = mr.applyEnvMappings(result, tenant)

		if len(result) == 0 {
			return nil, errMappingsFound
		}

		log.Debug("tenant mappings", "value", result)
		return result, nil
	case err != nil:
		log.Error("unknown database error when fetching mappings", err, "tenant", tenant)

		return nil, err
	}

	for _, mapping := range dbMappings {
		mapping.Source = restql.DatabaseSource
		result[mapping.ResourceName()] = mapping
	}

	result = mr.applyEnvMappings(result, tenant)

	if len(result) == 0 {
		return nil, errMappingsFound
	}

	return result, nil
}

func (mr MappingsReader) applyEnvMappings(result map[string]restql.Mapping, tenant string) map[string]restql.Mapping {
	tenantMappings, found := mr.env[tenant]
	if found {
		for k, v := range tenantMappings {
			result[k] = v
		}
	}

	return result
}

// ErrSetResourceMappingNotAllowed is returned when trying to write a resource mapping on a resource stored on local or env.
var ErrSetResourceMappingNotAllowed = errors.New("a resource mapping must have a source of type database in order to provide writing operations")

// MappingsWriter is the entity that maps resource name to URL.
type MappingsWriter struct {
	log   restql.Logger
	db    Database
	env   map[string]map[string]restql.Mapping
	local map[string]map[string]restql.Mapping
}

// NewMappingWriter creates an instance of MappingsWriter
func NewMappingWriter(log restql.Logger, env domain.EnvSource, local map[string]map[string]string, db Database) MappingsWriter {
	envMappings := getMappingsFromEnv(log, env)
	localMappings := make(map[string]map[string]restql.Mapping)
	for t, m := range local {
		localMappings[t] = parseMappingsFromLocal(log, m)
	}

	return MappingsWriter{log: log, env: envMappings, local: localMappings, db: db}
}

// Create makes a new mapping from the name to the URL on the given tenant
func (mw *MappingsWriter) Create(ctx context.Context, tenant string, resource string, url string) error {
	return mw.db.CreateMapping(ctx, tenant, resource, url)
}

// Update sets a URL to a resource name under the given tenant
func (mw *MappingsWriter) Update(ctx context.Context, tenant string, resource string, url string) error {
	if !mw.allowWrite(tenant, resource) {
		log := restql.GetLogger(ctx)
		log.Error("write operation on resource mapping not allowed", ErrSetResourceMappingNotAllowed, "tenant", tenant, "resource", resource)
		return ErrSetResourceMappingNotAllowed
	}

	return mw.db.SetMapping(ctx, tenant, resource, url)
}

func (mw *MappingsWriter) allowWrite(tenant string, resourceName string) bool {
	localTenantMappings, found := mw.local[tenant]
	if found {
		for resource := range localTenantMappings {
			if resource == resourceName {
				return false
			}
		}
	}

	envTenantMappings, found := mw.env[tenant]
	if found {
		for resource := range envTenantMappings {
			if resource == resourceName {
				return false
			}
		}
	}

	return true
}

func getMappingsFromEnv(log restql.Logger, envSource domain.EnvSource) map[string]map[string]restql.Mapping {
	result := make(map[string]map[string]restql.Mapping)
	env := envSource.GetAll()

	for key, value := range env {
		matches := envMappingWithTenantRegex.FindAllStringSubmatch(key, -1)
		if len(matches) > 0 && len(matches[0]) >= 3 {
			tenant := matches[0][1]
			resource := strings.ToLower(matches[0][2])
			mapping, err := restql.NewMapping(resource, value)
			if err != nil {
				log.Error("failed to create mapping", err)
				continue
			}

			tenantMappings, found := result[tenant]
			if !found {
				result[tenant] = map[string]restql.Mapping{
					resource: mapping,
				}
				continue
			}

			mapping.Source = restql.EnvSource
			tenantMappings[resource] = mapping
		}
	}

	return result
}

func parseMappingsFromLocal(log restql.Logger, local map[string]string) map[string]restql.Mapping {
	result := make(map[string]restql.Mapping)
	for k, v := range local {
		mapping, err := restql.NewMapping(k, v)
		if err != nil {
			log.Error("failed to create mapping", err)
			continue
		}

		mapping.Source = restql.ConfigFileSource
		result[k] = mapping
	}

	return result
}
