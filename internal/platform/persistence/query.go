package persistence

import (
	"context"

	"github.com/b2wdigital/restQL-golang/v4/internal/domain"
	"github.com/b2wdigital/restQL-golang/v4/internal/eval"
	"github.com/b2wdigital/restQL-golang/v4/internal/platform/logger"
	"github.com/b2wdigital/restQL-golang/v4/internal/platform/persistence/database"
	"github.com/b2wdigital/restQL-golang/v4/pkg/restql"
	"github.com/pkg/errors"
)

type savedQueries map[string][]string

type QueryReader struct {
	log   *logger.Logger
	local map[string]savedQueries
	db    database.Database
}

func NewQueryReader(log *logger.Logger, local map[string]map[string][]string, db database.Database) QueryReader {
	l := make(map[string]savedQueries)
	for k, v := range local {
		l[k] = v
	}
	return QueryReader{log: log, local: l, db: db}
}

func (qr QueryReader) Get(ctx context.Context, namespace, id string, revision int) (domain.SavedQuery, error) {
	log := restql.GetLogger(ctx)

	savedQuery, err := qr.db.FindQuery(ctx, namespace, id, revision)
	if err != nil && err != database.ErrNoDatabase {
		log.Error("query not found in database", err, "namespace", namespace, "name", id, "revision", revision)
	}

	if savedQuery.Text != "" {
		return savedQuery, nil
	}

	localQuery, err := qr.getQueryFromLocal(namespace, id, revision)
	if err != nil {
		log.Info("query not found in local", "namespace", namespace, "name", id, "revision", revision)
		return domain.SavedQuery{}, eval.NotFoundError{Err: errors.Errorf("query not found: %s/%s/%d", namespace, id, revision)}
	}

	return domain.SavedQuery{Text: localQuery}, nil
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
