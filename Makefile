
dev:
	PORT=9000 HEALTH_PORT=9001 DEBUG_PORT=9002 ENV=development go run cmd/restQL/main.go

build:
	go build -o bin/restQL -ldflags="-s -w" cmd/restQL/main.go
