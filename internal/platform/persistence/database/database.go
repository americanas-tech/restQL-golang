package database

import (
	"context"
	"errors"
	"time"

	"github.com/b2wdigital/restQL-golang/v4/internal/domain"
	"github.com/b2wdigital/restQL-golang/v4/internal/platform/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type dbOptions struct {
	ConnectionTimeout time.Duration
	MappingsTimeout   time.Duration
	QueryTimeout      time.Duration
	DatabaseName      string
}

type Option func(o *dbOptions)

func WithConnectionTimeout(timeout time.Duration) Option {
	return func(o *dbOptions) {
		o.ConnectionTimeout = timeout
	}
}

func WithMappingsTimeout(timeout time.Duration) Option {
	return func(o *dbOptions) {
		o.MappingsTimeout = timeout
	}
}

func WithQueryTimeout(timeout time.Duration) Option {
	return func(o *dbOptions) {
		o.QueryTimeout = timeout
	}
}

func WithDatabaseName(name string) Option {
	return func(o *dbOptions) {
		o.DatabaseName = name
	}
}

type Database interface {
	FindMappingsForTenant(ctx context.Context, tenantId string) ([]domain.Mapping, error)
	FindQuery(ctx context.Context, namespace string, name string, revision int) (domain.SavedQuery, error)
}

func New(log *logger.Logger, connectionString string, optionList ...Option) (Database, error) {
	if connectionString == "" {
		log.Info("no database configuration detected")
		return noOpDatabase{}, nil
	}

	dbOptions := dbOptions{}
	for _, o := range optionList {
		o(&dbOptions)
	}

	timeout := dbOptions.ConnectionTimeout

	log.Info("starting database connection", "timeout", timeout.String())

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	client, err := mongo.Connect(ctx,
		options.Client().ApplyURI(connectionString),
		options.Client().SetConnectTimeout(timeout),
	)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	log.Info("database connection established", "url", connectionString)

	return mongoDatabase{logger: log, client: client, options: dbOptions}, nil
}

var ErrNoDatabase = errors.New("no op database")

type noOpDatabase struct{}

func (n noOpDatabase) FindMappingsForTenant(ctx context.Context, tenantId string) ([]domain.Mapping, error) {
	return []domain.Mapping{}, ErrNoDatabase
}

func (n noOpDatabase) FindQuery(ctx context.Context, namespace string, name string, revision int) (domain.SavedQuery, error) {
	return domain.SavedQuery{}, ErrNoDatabase
}
