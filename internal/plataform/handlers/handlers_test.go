package handlers_test

import (
	"context"
	"github.com/b2wdigital/restQL-golang/internal/plataform/handlers"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
	"net"
	"net/http"
	"strings"
	"testing"
)

func TestValidateQuery(t *testing.T) {
	body := strings.NewReader("from cart")
	request, _ := http.NewRequest(http.MethodPost, "http://localhost/validate-query", body)

	h := setupHandler(handlers.New())
	defer h.Close()

	response := h.ServeHttp(request)

	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, response.StatusCode)
	}
}

func BenchmarkValidateQuery(b *testing.B) {
	h := setupHandler(handlers.New())
	defer h.Close()

	b.Run("use fasthttp router", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			body := strings.NewReader("from cart")
			request, _ := http.NewRequest(http.MethodPost, "http://localhost/validate-query", body)

			response := h.ServeHttp(request)

			if response.StatusCode != http.StatusOK {
				b.Fatalf("expected %d, got %d", http.StatusOK, response.StatusCode)
			}
		}
	})
}

func setupHandler(h fasthttp.RequestHandler) *testHandler {
	ln := fasthttputil.NewInmemoryListener()
	go func() {
		err := fasthttp.Serve(ln, h)
		if err != nil {
			panic("failed to serve handler")
		}
	}()

	client := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return ln.Dial()
			},
		},
	}

	return &testHandler{
		ln:     ln,
		client: client,
	}
}

type testHandler struct {
	ln     *fasthttputil.InmemoryListener
	client *http.Client
}

func (h *testHandler) ServeHttp(r *http.Request) *http.Response {
	response, err := h.client.Do(r)
	if err != nil {
		panic(err)
	}

	return response
}

func (h *testHandler) Close() {
	h.ln.Close()
}
