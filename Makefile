
vcs_ref := $(shell git rev-parse HEAD)

dev:
	RESTQL_PORT=9000 RESTQL_HEALTH_PORT=9001 RESTQL_DEBUG_PORT=9002 RESTQL_ENV=development go run -ldflags="-X main.build=$(vcs_ref)" -mod=vendor cmd/restQL/main.go

e2e: e2e-up e2e-run

e2e-up:
	RESTQL_CONFIG=./test/e2e/restql.yml make dev

e2e-run:
	cd test/e2e && go test -count=1 ./...

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod vendor -o bin/restQL -ldflags="-s -w -X main.build=$(vcs_ref) -extldflags -static" -tags netgo cmd/restQL/main.go
