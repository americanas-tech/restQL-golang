package httpclient

import (
	"context"
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"github.com/b2wdigital/restQL-golang/internal/platform/conf"
	"github.com/b2wdigital/restQL-golang/internal/platform/logger"
	"github.com/pkg/errors"
	"github.com/rs/dnscache"
	"github.com/valyala/fasthttp"
	"net"
	"time"
)

type HttpClient struct {
	client *fasthttp.Client
	log    *logger.Logger
}

const network = "tcp"

func New(log *logger.Logger, cfg *conf.Config) HttpClient {
	clientCfg := cfg.Web.Client

	r := &dnscache.Resolver{}
	go func() {
		t := time.NewTicker(1 * time.Minute)
		defer t.Stop()
		for range t.C {
			r.Refresh(true)
		}
	}()
	c := &fasthttp.Client{
		Name:                     "restql",
		NoDefaultUserAgentHeader: false,
		ReadTimeout:              clientCfg.ReadTimeout,
		WriteTimeout:             clientCfg.WriteTimeout,
		MaxConnsPerHost:          clientCfg.MaxIdleConnectionsPerHosts,
		MaxIdleConnDuration:      clientCfg.MaxIdleConnDuration,
		MaxConnDuration:          clientCfg.MaxConnDuration,
		Dial: func(addr string) (conn net.Conn, err error) {
			host, port, err := net.SplitHostPort(addr)
			if err != nil {
				return nil, err
			}
			ips, err := r.LookupHost(context.Background(), host)
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
	}

	return HttpClient{client: c, log: log}
}

func (hc HttpClient) Do(ctx context.Context, request domain.HttpRequest) (domain.HttpResponse, error) {
	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(res)
	}()

	setupRequest(request, req)

	duration, err := hc.executeWithContext(ctx, req, res)
	switch {
	case err == domain.ErrRequestTimeout:
		hc.log.Debug("request execution did not complete on time", "request", request)
		return domain.HttpResponse{}, err
	case err != nil:
		return domain.HttpResponse{}, errors.Wrap(err, "request execution failed")
	}

	response, err := makeResponse(req, res, duration)
	if err != nil {
		return domain.HttpResponse{}, err
	}

	return response, nil
}

func (hc HttpClient) executeWithContext(ctx context.Context, req *fasthttp.Request, res *fasthttp.Response) (time.Duration, error) {
	var start time.Time

	errCh := make(chan error)
	go func() {
		start = time.Now()
		errCh <- hc.client.Do(req, res)
	}()

	select {
	case e := <-errCh:
		finish := time.Since(start)
		return finish, e
	case <-ctx.Done():
		finish := time.Since(start)
		return finish, domain.ErrRequestTimeout
	}
}
