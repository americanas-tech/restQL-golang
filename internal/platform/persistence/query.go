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
	return QueryReader{log: log, local: parseLocalQueries(local), db: db}
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

		return restql.SavedQuery{}, restql.ErrQueryNotFound
	case err != nil:
		log.Error("database error when fetching query", err, "namespace", namespace, "name", id, "revision", revision)
		if localQuery.Text != "" {
			return localQuery, nil
		}

		return restql.SavedQuery{}, err
	}

	if dbQuery.Text != "" {
		dbQuery.Source = restql.DatabaseSource
		return dbQuery, nil
	}

	if localQuery.Text != "" {
		return localQuery, nil
	}

	return restql.SavedQuery{}, restql.ErrQueryNotFound
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
		log := restql.GetLogger(ctx)
		log.Error("fail to find namespaces on database", err)
	} else {
		for _, namespace := range dbNamespaces {
			namespaceSet[namespace] = struct{}{}
		}
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
		log := restql.GetLogger(ctx)
		log.Error("fail to find queries for namespace on database", err)
	}

	if len(dbQueries) == 0 && len(queries) == 0 {
		return nil, restql.ErrNamespaceNotFound
	}

	for queryName, dbRevisions := range dbQueries {
		dbRevisions = setDatabaseSource(dbRevisions)
		localRevisions := queries[queryName]
		queries[queryName] = unionLists(localRevisions, dbRevisions)
	}

	return queries, nil
}

func (qr QueryReader) ListQueryRevisions(ctx context.Context, namespace string, queryName string) ([]restql.SavedQuery, error) {
	var localRevisions []restql.SavedQuery

	namespaceQueries, found := qr.local[namespace]
	if found {
		revisions, found := namespaceQueries[queryName]
		if found {
			localRevisions = revisions
		}
	}

	dbRevisions, err := qr.db.FindQueryWithAllRevisions(ctx, namespace, queryName)
	if err != nil {
		log := restql.GetLogger(ctx)
		log.Error("fail to find query revisions from database", err)
	}

	dbRevisions = setDatabaseSource(dbRevisions)

	if len(localRevisions) == 0 && len(dbRevisions) == 0 {
		return nil, restql.ErrQueryNotFound
	}

	return unionLists(localRevisions, dbRevisions), nil
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

func setDatabaseSource(r []restql.SavedQuery) []restql.SavedQuery {
	for i, savedQuery := range r {
		savedQuery.Source = restql.DatabaseSource
		r[i] = savedQuery
	}

	return r
}

var ErrCreateRevisionNotAllowed = errors.New("a local query cannot be updated with a new revision : remove it from the local configuration or migrate it to the database")

type QueryWriter struct {
	log   restql.Logger
	local map[string]map[string][]restql.SavedQuery
	db    Database
}

func NewQueryWriter(log restql.Logger, local map[string]map[string][]string, db Database) QueryWriter {
	return QueryWriter{
		log:   log,
		local: parseLocalQueries(local),
		db:    db,
	}
}

func (qw QueryWriter) Write(ctx context.Context, namespace, name, content string) error {
	if !qw.allowWrite(ctx, namespace, name) {
		return ErrCreateRevisionNotAllowed
	}

	return qw.db.CreateQueryRevision(ctx, namespace, name, content)
}

func (qw QueryWriter) allowWrite(ctx context.Context, namespace string, name string) bool {
	namespaceQueries, found := qw.local[namespace]
	if !found {
		return true
	}

	queryRevisions, found := namespaceQueries[name]
	if !found {
		return true
	}

	dbRevisions, err := qw.db.FindQueryWithAllRevisions(ctx, namespace, name)
	if err != nil && err != restql.ErrQueryNotFoundInDatabase {
		return false
	}

	return len(dbRevisions) > 0 || len(queryRevisions) == 0
}

func parseLocalQueries(local map[string]map[string][]string) map[string]map[string][]restql.SavedQuery {
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
					Source:   restql.ConfigFileSource,
				}
			}

			parsedQueries[queryName] = parsedRevisions
		}

		l[namespace] = parsedQueries
	}
	return l
}
