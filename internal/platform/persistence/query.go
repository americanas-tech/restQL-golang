package persistence

import (
	"context"
	"github.com/b2wdigital/restQL-golang/v4/internal/eval"
	"github.com/b2wdigital/restQL-golang/v4/pkg/restql"
	"github.com/pkg/errors"
)

type savedQueries map[string][]string

// A QueryReader get a query from local configuration file
// or a database instance.
type QueryReader struct {
	log   restql.Logger
	local map[string]savedQueries
	db    Database
}

// NewQueryReader constructs a QueryReader from the given
// configuration and database.
func NewQueryReader(log restql.Logger, local map[string]map[string][]string, db Database) QueryReader {
	l := make(map[string]savedQueries)
	for k, v := range local {
		l[k] = v
	}
	return QueryReader{log: log, local: l, db: db}
}

// Get retrieves a query by its identity (namespace, id and revision),
// it first search the database and, if not found, in the configuration file.
func (qr QueryReader) Get(ctx context.Context, namespace, id string, revision int) (restql.SavedQuery, error) {
	log := restql.GetLogger(ctx)

	savedQuery, err := qr.db.FindQuery(ctx, namespace, id, revision)
	if err != nil && err != errNoDatabase {
		log.Error("query not found in database", err, "namespace", namespace, "name", id, "revision", revision)
	}

	if savedQuery.Text != "" {
		return savedQuery, nil
	}

	localQuery, err := qr.getQueryFromLocal(namespace, id, revision)
	if err != nil {
		log.Info("query not found in local", "namespace", namespace, "name", id, "revision", revision)
		return restql.SavedQuery{}, eval.NotFoundError{Err: errors.Errorf("query not found: %s/%s/%d", namespace, id, revision)}
	}

	return restql.SavedQuery{Text: localQuery}, nil
}

func (qr QueryReader) getQueryFromLocal(namespace string, id string, revision int) (string, error) {
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
