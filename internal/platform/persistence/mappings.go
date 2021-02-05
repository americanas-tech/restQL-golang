package persistence

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/b2wdigital/restQL-golang/v4/internal/domain"
	"github.com/b2wdigital/restQL-golang/v4/pkg/restql"
)

var envMappingRegex = regexp.MustCompile("^RESTQL_MAPPING_([^_]+)$")
var envMappingWithTenantRegex = regexp.MustCompile("^RESTQL_MAPPING_(\\w+)_(\\w+)$")

// MappingsReader fetch indexed mappings from database,
// configuration file and environment variable.
type MappingsReader struct {
	log           restql.Logger
	env           map[string]restql.Mapping
	envWithTenant map[string]map[string]restql.Mapping
	local         map[string]restql.Mapping
	localByTenant map[string]map[string]restql.Mapping
	db            Database
}

// NewMappingReader constructs a MappingsReader instance.
func NewMappingReader(log restql.Logger, env domain.EnvSource, local map[string]string, localByTenant map[string]map[string]string, db Database) MappingsReader {
	envMappings := getMappingsFromEnv(log, env)
	envWithTenantMappings := getMappingsFromEnvWithTenant(log, env)
	localMappings := parseMappingsFromLocal(log, local)
	localMappingsByTenant := make(map[string]map[string]restql.Mapping)
	for t, m := range localByTenant {
		localMappingsByTenant[t] = parseMappingsFromLocal(log, m)
	}

	return MappingsReader{log: log, env: envMappings, envWithTenant: envWithTenantMappings, local: localMappings, localByTenant: localMappingsByTenant, db: db}
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

	for tenant := range mr.envWithTenant {
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

	for k, v := range mr.local {
		result[k] = v
	}

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

	log.Debug("tenant mappings", "value", result)
	return result, nil
}

func (mr MappingsReader) applyEnvMappings(result map[string]restql.Mapping, tenant string) map[string]restql.Mapping {
	for k, v := range mr.env {
		result[k] = v
	}

	tenantMappings, found := mr.envWithTenant[tenant]
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
	log           restql.Logger
	db            Database
	env           map[string]restql.Mapping
	envWithTenant map[string]map[string]restql.Mapping
	local         map[string]restql.Mapping
	localByTenant map[string]map[string]restql.Mapping
}

// NewMappingWriter creates an instance of MappingsWriter
func NewMappingWriter(log restql.Logger, env domain.EnvSource, local map[string]string, localByTenant map[string]map[string]string, db Database) MappingsWriter {
	envMappings := getMappingsFromEnv(log, env)
	envWithTenantMappings := getMappingsFromEnvWithTenant(log, env)
	localMappings := parseMappingsFromLocal(log, local)
	localMappingsByTenant := make(map[string]map[string]restql.Mapping)
	for t, m := range localByTenant {
		localMappingsByTenant[t] = parseMappingsFromLocal(log, m)
	}

	return MappingsWriter{log: log, env: envMappings, envWithTenant: envWithTenantMappings, local: localMappings, localByTenant: localMappingsByTenant, db: db}
}

// Write sets a URL to a resource name under the given tenant
func (mw *MappingsWriter) Write(ctx context.Context, tenant string, resource string, url string) error {
	if !mw.allowWrite(tenant, resource) {
		log := restql.GetLogger(ctx)
		log.Error("write operation on resource mapping not allowed", ErrSetResourceMappingNotAllowed, "tenant", tenant, "resource", resource)
		return ErrSetResourceMappingNotAllowed
	}

	return mw.db.SetMapping(ctx, tenant, resource, url)
}

func (mw *MappingsWriter) allowWrite(tenant string, resourceName string) bool {
	for resource := range mw.local {
		if resource == resourceName {
			return false
		}
	}

	for resource := range mw.env {
		if resource == resourceName {
			return false
		}
	}

	localTenantMappings, found := mw.localByTenant[tenant]
	if found {
		for resource := range localTenantMappings {
			if resource == resourceName {
				return false
			}
		}
	}

	envTenantMappings, found := mw.envWithTenant[tenant]
	if found {
		for resource := range envTenantMappings {
			if resource == resourceName {
				return false
			}
		}
	}

	return true
}

func getMappingsFromEnv(log restql.Logger, envSource domain.EnvSource) map[string]restql.Mapping {
	result := make(map[string]restql.Mapping)
	env := envSource.GetAll()

	for key, value := range env {
		matches := envMappingRegex.FindAllStringSubmatch(key, -1)
		if len(matches) > 0 && len(matches[0]) >= 2 {
			resource := strings.ToLower(matches[0][1])
			mapping, err := restql.NewMapping(resource, value)
			if err != nil {
				log.Error("failed to create mapping", err)
				continue
			}

			mapping.Source = restql.EnvSource
			result[resource] = mapping
		}
	}

	return result
}

func getMappingsFromEnvWithTenant(log restql.Logger, envSource domain.EnvSource) map[string]map[string]restql.Mapping {
	result := make(map[string]map[string]restql.Mapping)
	env := envSource.GetAll()

	for key, value := range env {
		matches := envMappingWithTenantRegex.FindAllStringSubmatch(key, -1)
		if len(matches) > 0 && len(matches[0]) >= 3 {
			tenant := strings.ToLower(matches[0][1])
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
