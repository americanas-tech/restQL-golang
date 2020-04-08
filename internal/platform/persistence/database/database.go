package database

import (
	"context"
	"errors"
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"github.com/b2wdigital/restQL-golang/internal/platform/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
	"time"
)

type dbOptions struct {
	ConnectionTimeout string
	MappingsTimeout   time.Duration
	QueryTimeout      time.Duration
}

type Option func(o *dbOptions)

func WithConnectionTimeout(timeout string) Option {
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

type Database interface {
	FindMappingsForTenant(ctx context.Context, tenantId string) ([]domain.Mapping, error)
	FindQuery(ctx context.Context, namespace string, name string, revision int) (string, error)
}

func New(log *logger.Logger, connectionString string, optionList ...Option) (Database, error) {
	if connectionString == "" {
		log.Info("no database configuration detected")
		return noOpDatabase{}, nil
	}

	//todo: fazer uso de timeouts para mappings e query
	dbOptions := dbOptions{}
	for _, o := range optionList {
		o(&dbOptions)
	}

	timeout, err := parseTimeout(dbOptions.ConnectionTimeout)
	if err != nil {
		return noOpDatabase{}, err
	}

	log.Info("starting database connection", "timeout-in-ms", timeout)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	client, err := mongo.Connect(ctx,
		options.Client().ApplyURI(connectionString),
		options.Client().SetConnectTimeout(timeout),
	)
	if err != nil {
		return noOpDatabase{}, err
	}

	log.Info("database connection established", "url", connectionString)

	return mongoDatabase{logger: log, client: client}, nil
}

func parseTimeout(timeout string) (time.Duration, error) {
	n, err := strconv.Atoi(timeout)
	if err != nil {
		return 0, err
	}

	return time.Millisecond * time.Duration(n), nil
}

var ErrNoDatabase = errors.New("no op database")

type noOpDatabase struct{}

func (n noOpDatabase) FindMappingsForTenant(ctx context.Context, tenantId string) ([]domain.Mapping, error) {
	return []domain.Mapping{}, ErrNoDatabase
}

func (n noOpDatabase) FindQuery(ctx context.Context, namespace string, name string, revision int) (string, error) {
	return "", ErrNoDatabase
}
