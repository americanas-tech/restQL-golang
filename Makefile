
vcs_ref := $(shell git rev-parse HEAD)

dev:
	PORT=9000 HEALTH_PORT=9001 DEBUG_PORT=9002 ENV=development go run -ldflags="-X main.build=$(vcs_ref)" cmd/restQL/main.go

build:
	go build -o bin/restQL -ldflags="-s -w -X main.build=$(vcs_ref)" cmd/restQL/main.go
