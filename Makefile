
dev:
	PORT=9000 HEALTH_PORT=9001 go run cmd/restQL/main.go

build:
	go build -o bin/restQL -ldflags="-s -w" cmd/restQL/main.go
