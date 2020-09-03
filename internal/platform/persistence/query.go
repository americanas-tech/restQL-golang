package persistence

import (
	"context"
	"fmt"
	"github.com/b2wdigital/restQL-golang/v4/internal/domain"
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
	queryNotFoundErr := fmt.Errorf("%w: %s/%s/%d", domain.ErrQueryNotFound, namespace, id, revision)

	localQueryText, err := qr.getQueryFromLocal(namespace, id, revision)
	if err != nil {
		log.Info("query not found in local", "error", err, "namespace", namespace, "name", id, "revision", revision)
	}
	localQuery := restql.SavedQuery{Text: localQueryText}

	dbQuery, err := qr.db.FindQuery(ctx, namespace, id, revision)
	switch {
	case errors.Is(err, restql.ErrQueryNotFoundInDatabase):
		log.Error("query not found in database", err, "namespace", namespace, "name", id, "revision", revision)
		if localQuery.Text != "" {
			return localQuery, nil
		}

		return restql.SavedQuery{}, queryNotFoundErr
	case errors.Is(err, restql.ErrDatabaseCommunicationFailed):
		log.Error("database communication failed when fetching query", err, "namespace", namespace, "name", id, "revision", revision)
		if localQuery.Text != "" {
			return localQuery, nil
		}

		return restql.SavedQuery{}, err
	case err == errNoDatabase:
		if localQuery.Text != "" {
			return localQuery, nil
		}

		return restql.SavedQuery{}, queryNotFoundErr
	case err != nil:
		log.Error("unknown database error when fetching query", err, "namespace", namespace, "name", id, "revision", revision)
		if localQuery.Text != "" {
			return localQuery, nil
		}

		return restql.SavedQuery{}, err
	}

	if dbQuery.Text != "" {
		return dbQuery, nil
	}

	if localQuery.Text != "" {
		return localQuery, nil
	}

	return restql.SavedQuery{}, queryNotFoundErr
}

func (qr QueryReader) getQueryFromLocal(namespace string, id string, revision int) (string, error) {
	queriesInNamespace, ok := qr.local[namespace]
	if !ok {
		return "", errors.Errorf("namespace not found in local: %s", namespace)
	}

	queriesByRevision, ok := queriesInNamespace[id]
	if !ok {
		return "", errors.Errorf("query not found in local: %s", id)
	}

	if len(queriesByRevision) < revision {
		return "", errors.Errorf("revision not found in local: %d", revision)
	}

	queryTxt := queriesByRevision[revision-1]
	return queryTxt, nil
}
