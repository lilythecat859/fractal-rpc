.PHONY: dev build test lint clean

dev:
	docker compose -f scripts/docker-compose.dev.yml up -d
build:
	CGO_ENABLED=1 go build -ldflags="-s -w" -o fractal ./cmd/fractal
test:
	go test -race ./...
lint:
	golangci-lint run ./...
clean:
	rm -f fractal
