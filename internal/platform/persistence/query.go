package persistence

import (
	"context"
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

	localQueryText, err := qr.getQueryFromLocal(namespace, id, revision)
	if err != nil {
		log.Info("query not found in local", "error", err, "namespace", namespace, "name", id, "revision", revision)
	}
	localQuery := restql.SavedQuery{Name: id, Text: localQueryText, Revision: revision}

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
		rs := make([]restql.SavedQuery, len(revisions))
		for i, text := range revisions {
			rs[i] = restql.SavedQuery{
				Name:     queryName,
				Text:     text,
				Revision: i + 1,
			}
		}

		queries[queryName] = rs
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
