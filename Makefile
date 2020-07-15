
vcs_ref := $(shell git rev-parse HEAD)

setup-peg:
	go get github.com/mna/pigeon@v1.0.0

peg:
	pigeon ./internal/parser/ast/grammar.peg > ./internal/parser/ast/grammar.go

dev:
	RESTQL_PORT=9000 RESTQL_HEALTH_PORT=9001 RESTQL_DEBUG_PORT=9002 RESTQL_ENV=development go run -race -ldflags="-X main.build=$(vcs_ref)" -mod=vendor cmd/restQL/main.go

unit:
	cd internal && go test -race -mod=vendor -count=1 ./...

e2e: e2e-up e2e-run

e2e-up:
	RESTQL_CONFIG=./test/e2e/restql.yml make dev

e2e-run:
	cd test/e2e && go test -count=1 ./...

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod vendor -o bin/restQL -ldflags="-s -w -X main.build=$(vcs_ref) -extldflags -static" -tags netgo cmd/restQL/main.go
