package restql

import (
	"context"
	"github.com/b2wdigital/restQL-golang/v4/internal/domain"
	"log"
	"net/http"
	"net/url"
	"sync"
)

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

const (
	LifecyclePluginType PluginType = iota
	DatabasePluginType
)

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

type PluginInfo struct {
	Name string
	Type PluginType
	New  func(Logger) (Plugin, error)
}

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

type LifecyclePlugin interface {
	Plugin
	BeforeTransaction(ctx context.Context, tr TransactionRequest) context.Context
	AfterTransaction(ctx context.Context, tr TransactionResponse) context.Context
	BeforeQuery(ctx context.Context, query string, queryCtx QueryContext) context.Context
	AfterQuery(ctx context.Context, query string, result map[string]interface{}) context.Context
	BeforeRequest(ctx context.Context, request HttpRequest) context.Context
	AfterRequest(ctx context.Context, request HttpRequest, response HttpResponse, err error) context.Context
}

type TransactionRequest struct {
	Url    *url.URL
	Method string
	Header http.Header
}

type TransactionResponse struct {
	Status int
	Header http.Header
	Body   []byte
}

type HttpRequest = domain.HttpRequest
type HttpResponse = domain.HttpResponse

type DatabasePlugin interface {
	Plugin
	FindMappingsForTenant(ctx context.Context, tenantId string) ([]Mapping, error)
	FindQuery(ctx context.Context, namespace string, name string, revision int) (SavedQuery, error)
}
