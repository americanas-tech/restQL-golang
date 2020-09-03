package restql

import (
	"context"
	"errors"
	"github.com/b2wdigital/restQL-golang/v4/internal/domain"
	"log"
	"net/http"
	"net/url"
	"sync"
)

// Plugin is the root interface that allows general
// handling of the plugins instance when need.
type Plugin interface {
	Name() string
}

var (
	plugins   pluginIndex
	pluginsMu sync.RWMutex
)

type pluginIndex struct {
	lifecycle []PluginInfo
	dbPlugin  *PluginInfo
}

// Plugin types
const (
	LifecyclePluginType PluginType = iota
	DatabasePluginType
)

// PluginType is an enum of possible plugin types supported by restQL,
// currently supports LifecyclePluginType and DatabasePluginType.
type PluginType int

func (pt PluginType) String() string {
	switch pt {
	case LifecyclePluginType:
		return "Lifecycle"
	case DatabasePluginType:
		return "Database"
	default:
		return "Unknown"
	}
}

// PluginInfo represents a plugin instance associating a
// name and type to a constructor function.
type PluginInfo struct {
	Name string
	Type PluginType
	New  func(Logger) (Plugin, error)
}

// RegisterPlugin indexes the provided plugin information
// for latter usage by restQL in runtime.
// It supports registration of multiple Lifecycle plugins
// but only one Database plugin.
// In case of failure to register the plugin a warn
// message will be printed to the os.Stdout.
func RegisterPlugin(pluginInfo PluginInfo) {
	pluginsMu.Lock()
	defer pluginsMu.Unlock()

	switch pluginInfo.Type {
	case LifecyclePluginType:
		plugins.lifecycle = append(plugins.lifecycle, pluginInfo)
	case DatabasePluginType:
		if plugins.dbPlugin != nil {
			log.Printf("[WARN] database plugin already registred: %s", plugins.dbPlugin.Name)
			return
		}

		plugins.dbPlugin = &pluginInfo
	default:
		log.Printf("[WARN] unknown plugin type: %s", pluginInfo.Type)
	}
}

func GetLifecyclePlugins() []PluginInfo {
	pluginsMu.RLock()
	defer pluginsMu.RUnlock()

	lp := plugins.lifecycle

	return lp
}

func GetDatabasePlugin() (PluginInfo, bool) {
	pluginsMu.RLock()
	defer pluginsMu.RUnlock()

	dbPlugin := plugins.dbPlugin
	if dbPlugin == nil {
		return PluginInfo{}, false
	}

	return *dbPlugin, true
}

// LifecyclePlugin is the interface that defines
// all possible hooks during the query execution.
type LifecyclePlugin interface {
	Plugin
	BeforeTransaction(ctx context.Context, tr TransactionRequest) context.Context
	AfterTransaction(ctx context.Context, tr TransactionResponse) context.Context
	BeforeQuery(ctx context.Context, query string, queryCtx QueryContext) context.Context
	AfterQuery(ctx context.Context, query string, result map[string]interface{}) context.Context
	BeforeRequest(ctx context.Context, request HttpRequest) context.Context
	AfterRequest(ctx context.Context, request HttpRequest, response HttpResponse, err error) context.Context
}

// TransactionRequest represents a query execution
// transaction received through the /run-query/* endpoints.
type TransactionRequest struct {
	Url    *url.URL
	Method string
	Header http.Header
}

// TransactionResponse represents a query execution result
// from a transaction received through the /run-query/* endpoints.
type TransactionResponse struct {
	Status int
	Header http.Header
	Body   []byte
}

// HttpRequest represents a HTTP call to be
// made to an upstream dependency defined
// by the mappings.
type HttpRequest = domain.HTTPRequest

// HttpResponse represents a HTTP call result
// from an upstream dependency defined
// by the mappings.
type HttpResponse = domain.HTTPResponse

// DatabasePlugin is the interface that defines
// the obligatory operations needed from a database.
type DatabasePlugin interface {
	Plugin
	FindMappingsForTenant(ctx context.Context, tenantID string) ([]Mapping, error)
	FindQuery(ctx context.Context, namespace string, name string, revision int) (SavedQuery, error)
}

// Errors returned by Database plugin
var (
	ErrMappingsNotFoundInDatabase  = errors.New("mappings not found in database")
	ErrQueryNotFoundInDatabase     = errors.New("query not found in database")
	ErrDatabaseCommunicationFailed = errors.New("failed to communicate with the database")
)
