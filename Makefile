setup-peg:
	go get github.com/mna/pigeon@v1.0.0

peg:
	pigeon ./internal/parser/ast/grammar.peg > ./internal/parser/ast/grammar.go

dev:
	RESTQL_PORT=9000 RESTQL_HEALTH_PORT=9001 RESTQL_PPROF_PORT=9002 go run -race -ldflags="-X github.com/b2wdigital/restQL-golang/v4/cmd.build=$(RESTQL_BUILD)" main.go

unit:
	go test -race -count=1 ./internal/...
	go test -race -count=1 ./pkg/...

e2e:
	make e2e-up &
	sleep 10
	make e2e-run

e2e-up:
	RESTQL_CONFIG=./test/e2e/restql.yml make dev

e2e-run:
	go test -race -count=1 ./test/e2e/...

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/restQL -ldflags="-s -w -X github.com/b2wdigital/restQL-golang/v4/cmd.build=$(RESTQL_BUILD) -extldflags -static" -tags netgo main.go

# Modules support
deps-reset:
	git checkout -- go.mod
	go mod tidy
	go mod vendor

tidy:
	go mod tidy
	go mod vendor

deps-upgrade:
	go get -u -t -d -v ./...
	go mod tidy
	go mod vendor

deps-cleancache:
	go clean -modcache
#	
