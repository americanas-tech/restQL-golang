package httpclient

import (
	"context"
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"github.com/b2wdigital/restQL-golang/internal/platform/conf"
	"github.com/rs/dnscache"
	"github.com/valyala/fasthttp"
	"net"
	"strings"
	"sync"
	"time"
)

type clientPool struct {
	cfg      *conf.Config
	dialer   *fasthttp.TCPDialer
	clients  map[string]*fasthttp.HostClient
	mClients sync.Mutex
}

func newClientPool(cfg *conf.Config) *clientPool {
	r := &dnscache.Resolver{}
	go func() {
		t := time.NewTicker(10 * time.Minute)
		defer t.Stop()
		for range t.C {
			r.Refresh(true)
		}
	}()
	dialer := &fasthttp.TCPDialer{
		Resolver: &net.Resolver{
			PreferGo:     true,
			StrictErrors: false,
			Dial: func(ctx context.Context, network, address string) (conn net.Conn, err error) {
				host, port, err := net.SplitHostPort(address)
				if err != nil {
					return nil, err
				}
				ips, err := r.LookupHost(ctx, host)
				if err != nil {
					return nil, err
				}
				for _, ip := range ips {
					var dialer net.Dialer
					conn, err = dialer.Dial(network, net.JoinHostPort(ip, port))
					if err == nil {
						break
					}
				}
				return
			},
		},
	}

	return &clientPool{
		cfg:     cfg,
		clients: make(map[string]*fasthttp.HostClient),
		dialer:  dialer,
	}
}

func (cp *clientPool) Get(request domain.HttpRequest) *fasthttp.HostClient {
	clientCfg := cp.cfg.Web.Client

	host := request.Host
	isTLS := request.Schema == "https"

	var client *fasthttp.HostClient

	client, found := cp.clients[host]
	if !found {
		client = &fasthttp.HostClient{
			Addr:                          addMissingPort(host, isTLS),
			Name:                          "restql-" + host,
			NoDefaultUserAgentHeader:      false,
			DisableHeaderNamesNormalizing: true,
			IsTLS:                         isTLS,
			Dial:                          cp.dialer.Dial,
			ReadTimeout:                   clientCfg.ReadTimeout,
			WriteTimeout:                  clientCfg.WriteTimeout,
			MaxConns:                      clientCfg.MaxConnsPerHost,
			MaxIdleConnDuration:           clientCfg.MaxIdleConnDuration,
			MaxConnDuration:               clientCfg.MaxConnDuration,
			MaxConnWaitTimeout:            clientCfg.MaxConnWaitTimeout,
		}

		cp.mClients.Lock()
		cp.clients[host] = client
		cp.mClients.Unlock()
	}

	return client
}

func addMissingPort(addr string, isTLS bool) string {
	n := strings.Index(addr, ":")
	if n >= 0 {
		return addr
	}
	port := "80"
	if isTLS {
		port = "443"
	}
	return net.JoinHostPort(addr, port)
}
