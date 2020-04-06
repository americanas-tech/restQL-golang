package persistence

import (
	"github.com/b2wdigital/restQL-golang/internal/eval"
	"github.com/pkg/errors"
)

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
