package persistence

import (
	"context"
	"regexp"
	"strings"

	"github.com/b2wdigital/restQL-golang/v4/internal/domain"
	"github.com/b2wdigital/restQL-golang/v4/internal/platform/logger"
	"github.com/b2wdigital/restQL-golang/v4/pkg/restql"
)

var envMappingRegex = regexp.MustCompile("^RESTQL_MAPPING_(\\w+)")

type MappingsReader struct {
	log   *logger.Logger
	env   map[string]restql.Mapping
	local map[string]restql.Mapping
	db    Database
}

func NewMappingReader(log *logger.Logger, env domain.EnvSource, local map[string]string, db Database) MappingsReader {
	envMappings := getMappingsFromEnv(log, env)
	localMappings := parseMappingsFromLocal(log, local)

	return MappingsReader{log: log, env: envMappings, local: localMappings, db: db}
}

func (mr MappingsReader) FromTenant(ctx context.Context, tenant string) (map[string]restql.Mapping, error) {
	log := restql.GetLogger(ctx)
	log.Debug("fetching mappings")

	result := make(map[string]restql.Mapping)

	for k, v := range mr.local {
		result[k] = v
	}

	dbMappings, err := mr.db.FindMappingsForTenant(ctx, tenant)
	if err != nil && err != ErrNoDatabase {
		log.Debug("failed to load mappings from database", "error", err)
		return nil, err
	}

	for _, mapping := range dbMappings {
		result[mapping.ResourceName()] = mapping
	}

	for k, v := range mr.env {
		result[k] = v
	}

	log.Debug("tenant mappings", "value", result)

	return result, nil
}

func getMappingsFromEnv(log *logger.Logger, envSource domain.EnvSource) map[string]restql.Mapping {
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

func parseMappingsFromLocal(log *logger.Logger, local map[string]string) map[string]restql.Mapping {
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
