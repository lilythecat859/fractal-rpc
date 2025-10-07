.PHONY: build test run docker clean release lint
BINARY := fractal
GOFILES := $(shell find . -name '*.go' -type f)

build: $(BINARY)

$(BINARY): go.mod $(GOFILES)
	CGO_ENABLED=1 go build -ldflags="-s -w" -o $@ ./cmd/fractal

test:
	go test -race ./...

lint:
	golangci-lint run ./...

run: build
	./$(BINARY) -config fractal.toml

docker:
	docker build -t ghcr.io/lilythecat859/fractal:latest .

clean:
	rm -f $(BINARY)

release: clean
	./scripts/release.sh $(VERSION)

install: build
	install -D -m 755 $(BINARY) /usr/local/bin/$(BINARY)
	install -D -m 644 fractal.toml.example /etc/fractal.toml.example
