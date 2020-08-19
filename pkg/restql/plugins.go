package restql

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"sync"

	"github.com/b2wdigital/restQL-golang/v4/internal/domain"
)

type Plugin interface {
	Name() string
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

type QueryInput = domain.QueryInput
type QueryOptions = domain.QueryOptions
type QueryContext = domain.QueryContext

type HttpRequest = domain.HttpRequest
type HttpResponse = domain.HttpResponse

var (
	plugins   pluginIndex
	pluginsMu sync.RWMutex
)

type pluginIndex struct {
	lifecycle []PluginInfo
}

const (
	Lifecycle PluginType = iota
	Database
)

type PluginType int

func (pt PluginType) String() string {
	switch pt {
	case Lifecycle:
		return "Lifecycle"
	case Database:
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

func RegisterPlugin(loader func() PluginInfo) {
	pluginsMu.Lock()
	defer pluginsMu.Unlock()

	pi := loader()

	switch pi.Type {
	case Lifecycle:
		plugins.lifecycle = append(plugins.lifecycle, pi)
	default:
		log.Printf("[WARN] unknown plugin type: %s", pi.Type)
	}
}

func LifecyclePlugins() []PluginInfo {
	pluginsMu.RLock()
	defer pluginsMu.RUnlock()

	lp := plugins.lifecycle

	return lp
}
