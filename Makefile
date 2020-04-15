
vcs_ref := $(shell git rev-parse HEAD)

dev:
	RESTQL_PORT=9000 RESTQL_HEALTH_PORT=9001 RESTQL_DEBUG_PORT=9002 RESTQL_ENV=development go run -ldflags="-X main.build=$(vcs_ref)" cmd/restQL/main.go

build:
	CGO_ENABLED=0 go build -mod vendor -o bin/restQL -ldflags="-s -w -X main.build=$(vcs_ref)" -tags netgo cmd/restQL/main.go
