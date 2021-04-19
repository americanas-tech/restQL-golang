package persistence

import (
	"context"

	"github.com/b2wdigital/restQL-golang/v6/pkg/restql"
	"github.com/pkg/errors"
)

// A QueryReader get a query from local configuration file
// or a database instance.
type QueryReader struct {
	log   restql.Logger
	local map[string][]restql.SavedQuery
	db    Database
}

// NewQueryReader constructs a QueryReader from the given
// configuration and database.
func NewQueryReader(log restql.Logger, local map[string]map[string][]string, db Database) QueryReader {
	return QueryReader{log: log, local: parseLocalQueries(local), db: db}
}

// Get retrieves a query by its identity (namespace, id and revision),
// it first search the database and, if not found, in the configuration file.
func (qr QueryReader) Get(ctx context.Context, namespace, id string, revision int) (restql.SavedQueryRevision, error) {
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

		return restql.SavedQueryRevision{}, restql.ErrQueryNotFound
	case err != nil:
		log.Error("database error when fetching query", err, "namespace", namespace, "name", id, "revision", revision)
		if localQuery.Text != "" {
			return localQuery, nil
		}

		return restql.SavedQueryRevision{}, err
	}

	if dbQuery.Text != "" {
		dbQuery.Source = restql.DatabaseSource
		return dbQuery, nil
	}

	if localQuery.Text != "" {
		return localQuery, nil
	}

	return restql.SavedQueryRevision{}, restql.ErrQueryNotFound
}

func (qr QueryReader) getQueryFromLocal(namespace string, id string, revision int) (restql.SavedQueryRevision, error) {
	queriesInNamespace, ok := qr.local[namespace]
	if !ok {
		return restql.SavedQueryRevision{}, errors.Errorf("namespace not found in local: %s", namespace)
	}

	var query restql.SavedQuery
	found := false
	for _, q := range queriesInNamespace {
		if q.Name == id {
			query = q
			found = true
		}
	}

	if !found {
		return restql.SavedQueryRevision{}, errors.Errorf("query not found in local: %s", id)
	}

	if len(query.Revisions) < revision {
		return restql.SavedQueryRevision{}, errors.Errorf("revision not found in local: %d", revision)
	}

	savedQuery := query.Revisions[revision-1]
	return savedQuery, nil
}

// ListNamespaces fetches all namespace on config, env and database
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

// ListQueriesForNamespace fetches the queries under the given namespace,
// stored on config file, env or database.
func (qr QueryReader) ListQueriesForNamespace(ctx context.Context, namespace string, archived bool) ([]restql.SavedQuery, error) {
	var queries []restql.SavedQuery
	if !archived {
		queries = qr.local[namespace]
	}

	dbQueries, err := qr.db.FindQueriesForNamespace(ctx, namespace, archived)
	if err != nil && err != restql.ErrNamespaceNotFound {
		log := restql.GetLogger(ctx)
		log.Error("fail to find queries for namespace on database", err)
	}

	if err == restql.ErrNamespaceNotFound && len(queries) == 0 {
		return nil, restql.ErrNamespaceNotFound
	}

	for _, dbQuery := range dbQueries {
		query := dbQuery
		query.Revisions = setDatabaseSource(dbQuery.Revisions)
		localQuery, found := findQueryByName(queries, query.Name)
		if found {
			query.Revisions = revisionsUnion(localQuery.Revisions, query.Revisions)
		}

		queries = append(queries, query)
	}

	qr.log.Debug("query reader namespace", "queries", queries)

	return queries, nil
}

// ListQueryRevisions fetch revisions for a query on the given namespace,
// stored on the config file, env or database.
func (qr QueryReader) ListQueryRevisions(ctx context.Context, namespace string, queryName string, archived bool) (restql.SavedQuery, error) {
	var localQuery restql.SavedQuery
	q, found := findQueryByName(qr.local[namespace], queryName)
	if found && !archived {
		localQuery = q
	}

	dbQuery, err := qr.db.FindQueryWithAllRevisions(ctx, namespace, queryName, archived)
	if err != nil && err != restql.ErrQueryNotFoundInDatabase {
		log := restql.GetLogger(ctx)
		log.Error("fail to find query revisions from database", err)
	}

	if err == restql.ErrQueryNotFoundInDatabase && len(localQuery.Revisions) == 0 {
		return restql.SavedQuery{}, restql.ErrQueryNotFound
	}

	query := dbQuery
	query.Revisions = setDatabaseSource(query.Revisions)
	query.Revisions = revisionsUnion(localQuery.Revisions, query.Revisions)

	return query, nil
}

func revisionsUnion(a, b []restql.SavedQueryRevision) []restql.SavedQueryRevision {
	if len(b) > len(a) {
		return b
	}

	result := make([]restql.SavedQueryRevision, len(a))
	copy(result, a)
	copy(result, b)
	return result
}

func setDatabaseSource(r []restql.SavedQueryRevision) []restql.SavedQueryRevision {
	for i, savedQuery := range r {
		savedQuery.Source = restql.DatabaseSource
		r[i] = savedQuery
	}

	return r
}

// ErrUpdateQueryNotAllowed is returned when trying to write a query revision
// on a query stored on local or env.
var ErrUpdateQueryNotAllowed = errors.New("a local query cannot be updated : remove it from the local configuration or migrate it to the database")

// QueryWriter is the entity that create a new query revision
// when it is stored on the database.
type QueryWriter struct {
	log   restql.Logger
	local map[string][]restql.SavedQuery
	db    Database
}

// NewQueryWriter creates a QueryWriter instance
func NewQueryWriter(log restql.Logger, local map[string]map[string][]string, db Database) QueryWriter {
	return QueryWriter{
		log:   log,
		local: parseLocalQueries(local),
		db:    db,
	}
}

// Write creates a new query revision
func (qw QueryWriter) Write(ctx context.Context, namespace, name, content string) error {
	if !qw.allowWrite(namespace, name) {
		return ErrUpdateQueryNotAllowed
	}

	return qw.db.CreateQueryRevision(ctx, namespace, name, content)
}

func (qw QueryWriter) UpdateQueryArchiving(ctx context.Context, namespace string, name string, archived bool) error {
	if !qw.allowWrite(namespace, name) {
		return ErrUpdateQueryNotAllowed
	}

	err := qw.db.UpdateQueryArchiving(ctx, namespace, name, archived)
	if err != nil {
		qw.log.Error("failed to update query archiving", err)
	}

	return err
}

func (qw QueryWriter) UpdateRevisionArchiving(ctx context.Context, namespace string, name string, revision int, archived bool) error {
	if !qw.allowWrite(namespace, name) {
		return ErrUpdateQueryNotAllowed
	}

	err := qw.db.UpdateRevisionArchiving(ctx, namespace, name, revision, archived)
	if err != nil {
		qw.log.Error("failed to update revision archiving", err)
	}

	return err
}

func (qw QueryWriter) allowWrite(namespace string, name string) bool {
	namespaceQueries, found := qw.local[namespace]
	if !found {
		return true
	}

	_, found = findQueryByName(namespaceQueries, name)
	return !found
}

func findQueryByName(queries []restql.SavedQuery, name string) (restql.SavedQuery, bool) {
	for _, query := range queries {
		if query.Name == name {
			return query, true
		}
	}

	return restql.SavedQuery{}, false
}

func parseLocalQueries(local map[string]map[string][]string) map[string][]restql.SavedQuery {
	l := make(map[string][]restql.SavedQuery)
	for namespace, queries := range local {
		namespacedQueries := make([]restql.SavedQuery, len(queries))

		j := 0
		for queryName, revisions := range queries {
			q := restql.SavedQuery{
				Namespace: namespace,
				Name:      queryName,
				Revisions: make([]restql.SavedQueryRevision, len(revisions)),
			}

			for i, text := range revisions {
				q.Revisions[i] = restql.SavedQueryRevision{
					Name:     queryName,
					Text:     text,
					Revision: i + 1,
					Source:   restql.ConfigFileSource,
				}
			}

			namespacedQueries[j] = q
			j++
		}

		l[namespace] = namespacedQueries
	}
	return l
}
