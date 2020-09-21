package persistence

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/b2wdigital/restQL-golang/v4/internal/domain"
	"github.com/b2wdigital/restQL-golang/v4/pkg/restql"
)

var envMappingRegex = regexp.MustCompile("^RESTQL_MAPPING_(\\w+)")

// MappingsReader fetch indexed mappings from database,
// configuration file and environment variable.
type MappingsReader struct {
	log   restql.Logger
	env   map[string]restql.Mapping
	local map[string]restql.Mapping
	db    Database
}

// NewMappingReader constructs a MappingsReader instance.
func NewMappingReader(log restql.Logger, env domain.EnvSource, local map[string]string, db Database) MappingsReader {
	envMappings := getMappingsFromEnv(log, env)
	localMappings := parseMappingsFromLocal(log, local)

	return MappingsReader{log: log, env: envMappings, local: localMappings, db: db}
}

// FromTenant fetch the mappings for the given tenant.
func (mr MappingsReader) FromTenant(ctx context.Context, tenant string) (map[string]restql.Mapping, error) {
	log := restql.GetLogger(ctx)
	log.Debug("fetching mappings")
	mappingsFoundErr := fmt.Errorf("%w: tenant %s", restql.ErrMappingsNotFoundInLocal, tenant)

	result := make(map[string]restql.Mapping)

	for k, v := range mr.local {
		result[k] = v
	}

	dbMappings, err := mr.db.FindMappingsForTenant(ctx, tenant)
	switch {
	case err == errNoDatabase:
		result = mr.applyEnvMappings(result)

		if len(result) == 0 {
			return nil, mappingsFoundErr
		}

		log.Debug("tenant mappings", "value", result)
		return result, nil
	case err != nil:
		log.Error("unknown database error when fetching mappings", err, "tenant", tenant)

		return nil, err
	}

	for _, mapping := range dbMappings {
		result[mapping.ResourceName()] = mapping
	}

	result = mr.applyEnvMappings(result)

	if len(result) == 0 {
		return nil, mappingsFoundErr
	}

	log.Debug("tenant mappings", "value", result)
	return result, nil
}

func (mr MappingsReader) applyEnvMappings(result map[string]restql.Mapping) map[string]restql.Mapping {
	for k, v := range mr.env {
		result[k] = v
	}
	return result
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

			result[resource] = mapping
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

		result[k] = mapping
	}

	return result
}
