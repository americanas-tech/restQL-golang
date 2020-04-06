package persistence

import (
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"regexp"
	"strings"
)

type MappingsReader struct {
	env   map[string]domain.Mapping
	local map[string]domain.Mapping
}

func NewMappingReader(env domain.EnvSource, local map[string]string) MappingsReader {
	envMappings := getMappingsFromEnv(env)
	localMappings := parseMappingsFromLocal(local)
	mr := MappingsReader{env: envMappings, local: localMappings}

	return mr
}

func parseMappingsFromLocal(local map[string]string) map[string]domain.Mapping {
	result := make(map[string]domain.Mapping)
	for k, v := range local {
		mapping, err := domain.NewMapping(k, v)
		if err != nil {
			//TODO: log error
			continue
		}

		result[k] = mapping
	}

	return result
}

//func (mr MappingsReader) Get(tenant, resource string) (domain.Mapping, error) {
//	switch {
//	case mr.env.GetString(resource) != "":
//		return domain.NewMapping(resource, mr.env.GetString(resource))
//	case mr.local[resource] != "":
//		return domain.NewMapping(resource, mr.local[resource])
//	default:
//		return domain.Mapping{}, eval.NotFoundError{Err: errors.Errorf("resource `%s` not found on mappings", resource)}
//	}
//}

func (mr MappingsReader) FromTenant(tenant string) (map[string]domain.Mapping, error) {
	result := make(map[string]domain.Mapping)

	for k, v := range mr.local {
		result[k] = v
	}

	for k, v := range mr.env {
		result[k] = v
	}

	return result, nil
}

var envMappingRegex = regexp.MustCompile("^RESTQL_MAPPING_(\\w+)")

func getMappingsFromEnv(envSource domain.EnvSource) map[string]domain.Mapping {
	result := make(map[string]domain.Mapping)
	env := envSource.GetAll()

	for key, value := range env {
		matches := envMappingRegex.FindAllStringSubmatch(key, -1)
		if len(matches) > 0 {
			resource := strings.ToLower(matches[0][1])
			mapping, err := domain.NewMapping(resource, value)
			if err != nil {
				//TODO: log error
				continue
			}

			result[resource] = mapping
		}
	}

	return result
}
