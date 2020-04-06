package persistence

import (
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"github.com/b2wdigital/restQL-golang/internal/eval"
	"github.com/b2wdigital/restQL-golang/internal/platform/conf"
	"github.com/pkg/errors"
)

type MappingsReader struct {
	env   conf.EnvSource
	local map[string]string
}

func NewMappingReader(env conf.EnvSource, local map[string]string) MappingsReader {
	mr := MappingsReader{env: env}
	mr.local = local

	return mr
}

func (mr MappingsReader) Get(tenant, resource string) (domain.Mapping, error) {
	switch {
	case mr.env.GetString(resource) != "":
		return domain.NewMapping(resource, mr.env.GetString(resource)), nil
	case mr.local[resource] != "":
		return domain.NewMapping(resource, mr.local[resource]), nil
	default:
		return domain.Mapping{}, eval.NotFoundError{Err: errors.Errorf("resource `%s` not found on mappings", resource)}
	}
}

type savedQueries map[string][]string

type QueryReader struct {
	local map[string]savedQueries
}

func NewQueryReader(local map[string]map[string][]string) QueryReader {
	l := make(map[string]savedQueries)
	for k, v := range local {
		l[k] = v
	}
	return QueryReader{local: l}
}

func (qr QueryReader) Get(namespace, id string, revision int) (string, error) {
	queriesInNamespace, ok := qr.local[namespace]
	if !ok {
		return "", eval.NotFoundError{Err: errors.Errorf("namespace not found: %s", namespace)}
	}

	queriesByRevision, ok := queriesInNamespace[id]
	if !ok {
		return "", eval.NotFoundError{Err: errors.Errorf("query not found: %s", id)}
	}

	if len(queriesByRevision) < revision {
		return "", eval.NotFoundError{Err: errors.Errorf("revision not found: %d", revision)}
	}

	queryTxt := queriesByRevision[revision-1]

	return queryTxt, nil
}
