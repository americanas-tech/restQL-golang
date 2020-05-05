package persistence

import (
	"context"
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"github.com/b2wdigital/restQL-golang/internal/platform/logger"
	"github.com/b2wdigital/restQL-golang/internal/platform/persistence/database"
	"regexp"
	"strings"
)

var envMappingRegex = regexp.MustCompile("^RESTQL_MAPPING_(\\w+)")

type MappingsReader struct {
	log   *logger.Logger
	env   map[string]domain.Mapping
	local map[string]domain.Mapping
	db    database.Database
}

func NewMappingReader(log *logger.Logger, env domain.EnvSource, local map[string]string, db database.Database) MappingsReader {
	envMappings := getMappingsFromEnv(log, env)
	localMappings := parseMappingsFromLocal(log, local)

	return MappingsReader{log: log, env: envMappings, local: localMappings, db: db}
}

func (mr MappingsReader) FromTenant(ctx context.Context, tenant string) (map[string]domain.Mapping, error) {
	mr.log.Debug("fetching mappings")

	result := make(map[string]domain.Mapping)

	for k, v := range mr.local {
		result[k] = v
	}

	dbMappings, err := mr.db.FindMappingsForTenant(ctx, tenant)
	if err != nil && err != database.ErrNoDatabase {
		mr.log.Debug("failed to load mappings from database", "error", err)
		return nil, err
	}

	for _, mapping := range dbMappings {
		result[mapping.ResourceName] = mapping
	}

	for k, v := range mr.env {
		result[k] = v
	}

	mr.log.Debug("tenant mappings", "value", result)

	return result, nil
}

func getMappingsFromEnv(log *logger.Logger, envSource domain.EnvSource) map[string]domain.Mapping {
	result := make(map[string]domain.Mapping)
	env := envSource.GetAll()

	for key, value := range env {
		matches := envMappingRegex.FindAllStringSubmatch(key, -1)
		if len(matches) > 0 {
			resource := strings.ToLower(matches[0][1])
			mapping, err := domain.NewMapping(resource, value)
			if err != nil {
				log.Error("failed to create mapping", err)
				continue
			}

			result[resource] = mapping
		}
	}

	return result
}

func parseMappingsFromLocal(log *logger.Logger, local map[string]string) map[string]domain.Mapping {
	result := make(map[string]domain.Mapping)
	for k, v := range local {
		mapping, err := domain.NewMapping(k, v)
		if err != nil {
			log.Error("failed to create mapping", err)
			continue
		}

		result[k] = mapping
	}

	return result
}
