package persistence

import (
	"context"
	"github.com/b2wdigital/restQL-golang/v4/pkg/restql"
	"github.com/pkg/errors"
)

// A QueryReader get a query from local configuration file
// or a database instance.
type QueryReader struct {
	log   restql.Logger
	local map[string]map[string][]restql.SavedQuery
	db    Database
}

// NewQueryReader constructs a QueryReader from the given
// configuration and database.
func NewQueryReader(log restql.Logger, local map[string]map[string][]string, db Database) QueryReader {
	l := make(map[string]map[string][]restql.SavedQuery)
	for namespace, queries := range local {
		parsedQueries := make(map[string][]restql.SavedQuery)
		for queryName, revisions := range queries {
			parsedRevisions := make([]restql.SavedQuery, len(revisions))
			for i, text := range revisions {
				parsedRevisions[i] = restql.SavedQuery{
					Name:     queryName,
					Text:     text,
					Revision: i + 1,
				}
			}

			parsedQueries[queryName] = parsedRevisions
		}

		l[namespace] = parsedQueries
	}
	return QueryReader{log: log, local: l, db: db}
}

// Get retrieves a query by its identity (namespace, id and revision),
// it first search the database and, if not found, in the configuration file.
func (qr QueryReader) Get(ctx context.Context, namespace, id string, revision int) (restql.SavedQuery, error) {
	log := restql.GetLogger(ctx)

	localQuery, err := qr.getQueryFromLocal(namespace, id, revision)
	if err != nil {
		log.Info("query not found in local", "error", err, "namespace", namespace, "name", id, "revision", revision)
	}

	dbQuery, err := qr.db.FindQuery(ctx, namespace, id, revision)
	switch {
	case err == errNoDatabase:
		if localQuery.Text != "" {
			return localQuery, nil
		}

		return restql.SavedQuery{}, restql.ErrQueryNotFoundInLocal
	case err != nil:
		log.Error("database error when fetching query", err, "namespace", namespace, "name", id, "revision", revision)
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

	return restql.SavedQuery{}, restql.ErrQueryNotFoundInDatabase
}

func (qr QueryReader) getQueryFromLocal(namespace string, id string, revision int) (restql.SavedQuery, error) {
	queriesInNamespace, ok := qr.local[namespace]
	if !ok {
		return restql.SavedQuery{}, errors.Errorf("namespace not found in local: %s", namespace)
	}

	queriesByRevision, ok := queriesInNamespace[id]
	if !ok {
		return restql.SavedQuery{}, errors.Errorf("query not found in local: %s", id)
	}

	if len(queriesByRevision) < revision {
		return restql.SavedQuery{}, errors.Errorf("revision not found in local: %d", revision)
	}

	savedQuery := queriesByRevision[revision-1]
	return savedQuery, nil
}

func (qr QueryReader) ListNamespaces(ctx context.Context) ([]string, error) {
	namespaceSet := make(map[string]struct{})

	for namespace := range qr.local {
		namespaceSet[namespace] = struct{}{}
	}

	dbNamespaces, err := qr.db.FindAllNamespaces(ctx)
	if err != nil {
		return nil, err
	}

	for _, namespace := range dbNamespaces {
		namespaceSet[namespace] = struct{}{}
	}

	namespaces := make([]string, len(namespaceSet))
	i := 0
	for namespace := range namespaceSet {
		namespaces[i] = namespace
		i++
	}

	return namespaces, nil
}

func (qr QueryReader) ListQueriesForNamespace(ctx context.Context, namespace string) (map[string][]restql.SavedQuery, error) {
	queries := make(map[string][]restql.SavedQuery)
	for queryName, revisions := range qr.local[namespace] {
		queries[queryName] = revisions
	}

	dbQueries, err := qr.db.FindQueriesForNamespace(ctx, namespace)
	if err != nil {
		return nil, err
	}

	for queryName, dbRevisions := range dbQueries {
		localRevisions := queries[queryName]

		if len(dbRevisions) > len(localRevisions) {
			queries[queryName] = dbRevisions
		} else {
			for i, dbr := range dbRevisions {
				localRevisions[i] = dbr
			}
			queries[queryName] = localRevisions
		}
	}

	return queries, nil
}

func (qr QueryReader) ListQueryRevisions(ctx context.Context, namespace string, queryName string) ([]restql.SavedQuery, error) {
	dbRevisions, err := qr.db.FindQueryWithAllRevisions(ctx, namespace, queryName)
	switch {
	case err == restql.ErrQueryNotFoundInDatabase:
		namespace, found := qr.local[namespace]
		if !found {
			return nil, restql.ErrQueryNotFoundInLocal
		}

		revisions, found := namespace[queryName]
		if !found {
			return nil, restql.ErrQueryNotFoundInLocal
		}

		return revisions, nil
	case err != nil:
		return nil, err
	default:
		namespace, found := qr.local[namespace]
		if !found {
			return dbRevisions, nil
		}

		localRevisions := namespace[queryName]

		return unionLists(localRevisions, dbRevisions), nil
	}
}

func unionLists(a, b []restql.SavedQuery) []restql.SavedQuery {
	if len(b) > len(a) {
		return b
	}

	result := make([]restql.SavedQuery, len(a))
	copy(result, a)

	for i, v := range b {
		result[i] = v
	}

	return result
}
