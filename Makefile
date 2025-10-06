.PHONY: build run test clean

BINARY=fractal
GO=go
CGO=1

build:
	CGO_ENABLED=$(CGO) $(GO) build -ldflags="-s -w" -o $(BINARY) ./cmd/fractal

run: build
	./$(BINARY)

test:
	$(GO) test -race ./...

clean:
	rm -f $(BINARY)

install: build
	sudo cp $(BINARY) /usr/local/bin/$(BINARY)
