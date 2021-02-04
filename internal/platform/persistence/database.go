package persistence

import (
	"context"
	"github.com/b2wdigital/restQL-golang/v4/pkg/restql"
	"github.com/pkg/errors"
)

// Database defines the operations exposed by an external store.
type Database restql.DatabasePlugin

// NewDatabase constructs a Database compliant value
// from the database plugin registered.
// In case of no plugin, a noop implementation is returned.
func NewDatabase(log restql.Logger, disabled bool) (Database, error) {
	if disabled {
		return noOpDatabase{}, nil
	}

	pluginInfo, found := restql.GetDatabasePlugin()
	if !found {
		log.Info("no database plugin provided")
		return noOpDatabase{}, nil
	}

	dbPlugin, err := pluginInfo.New(log)
	if err != nil {
		return noOpDatabase{}, err
	}

	if dbPlugin == nil {
		log.Info("empty database instance returned by plugin", "plugin", pluginInfo.Name)
		return noOpDatabase{}, nil
	}

	database, ok := dbPlugin.(restql.DatabasePlugin)
	if !ok {
		return noOpDatabase{}, errors.Errorf("failed to cast database plugin, unknown type: %T", dbPlugin)
	}

	return database, nil
}

var errNoDatabase = errors.New("no op database")

type noOpDatabase struct{}

func (n noOpDatabase) Name() string {
	return "noopdatabase"
}

func (n noOpDatabase) FindAllNamespaces(ctx context.Context) ([]string, error) {
	return []string{"demo", "cart", "checkout"}, nil

	//return nil, errNoDatabase
}

func (n noOpDatabase) FindQueriesForNamespace(ctx context.Context, namespace string) (map[string][]restql.SavedQuery, error) {
	return map[string][]restql.SavedQuery{
		"test": {
			{
				Name:     "test",
				Text:     "from test",
				Revision: 1,
			},
		},
	}, nil

	//return nil, errNoDatabase
}

func (n noOpDatabase) FindQueryWithAllRevisions(ctx context.Context, namespace string, queryName string) ([]restql.SavedQuery, error) {
	return []restql.SavedQuery{
		{
			Name:     "sku",
			Text:     "from test",
			Revision: 1,
		},
	}, nil

	//return nil, errNoDatabase
}

func (n noOpDatabase) CreateQueryRevision(ctx context.Context, namespace string, queryName string, content string) error {
	return errNoDatabase
}

func (n noOpDatabase) FindAllTenants(ctx context.Context) ([]string, error) {
	return []string{"dc", "marvel", "vertigo"}, nil

	//return nil, errNoDatabase
}

func (n noOpDatabase) SetMapping(ctx context.Context, tenantID string, mappingsName string, url string) error {
	return errNoDatabase
}

func (n noOpDatabase) FindMappingsForTenant(ctx context.Context, tenantID string) ([]restql.Mapping, error) {
	mapping, _ := restql.NewMapping("villain", "http://vilain.com")
	return []restql.Mapping{
		mapping,
	}, nil
}

func (n noOpDatabase) FindQuery(ctx context.Context, namespace string, name string, revision int) (restql.SavedQuery, error) {
	return restql.SavedQuery{}, errNoDatabase
}
