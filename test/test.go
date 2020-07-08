package test

import (
	"encoding/json"
	"fmt"
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"github.com/google/go-cmp/cmp"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"testing"
)

func Unmarshal(body string) interface{} {
	var f interface{}
	err := json.Unmarshal([]byte(body), &f)
	if err != nil {
		panic(err)
	}
	return f
}

var regexComparer = cmp.Comparer(func(x, y *regexp.Regexp) bool {
	return x.String() == y.String()
})

var mappingComparer = cmp.Comparer(func(x, y domain.Mapping) bool {
	return x.ResourceName() == y.ResourceName() &&
		x.Scheme() == y.Scheme() &&
		x.Host() == y.Host()
})

func Equal(t *testing.T, got, expected interface{}) {
	if !cmp.Equal(got, expected, regexComparer, mappingComparer) {
		t.Fatalf("got = %+#v, want = %+#v\nMismatch (-want +got):\n%s", got, expected, cmp.Diff(expected, got, regexComparer))
	}
}

func NotEqual(t *testing.T, got, expected interface{}) {
	if reflect.DeepEqual(got, expected) {
		t.Fatalf("got = %+#v, want = %+#v", got, expected)
	}
}

func VerifyError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("unexpected error returned: %s", err)
	}
}

type MockServer struct {
	port   int
	server *httptest.Server
	mux    *http.ServeMux
}

func NewMockServer(port int) *MockServer {
	mux := http.NewServeMux()
	mockServer := MockServer{port: port, mux: mux}

	return &mockServer
}

func (ms *MockServer) Start() {
	l, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", ms.port))
	if err != nil {
		log.Fatal(err)
	}

	server := httptest.NewUnstartedServer(ms.mux)
	server.Listener = l

	ms.server = server

	server.Start()
}

func (ms *MockServer) Mux() *http.ServeMux {
	return ms.mux
}

func (ms *MockServer) Server() *httptest.Server {
	return ms.server
}

func (ms *MockServer) Teardown() {
	ms.server.Close()
}
