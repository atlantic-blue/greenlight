BINARY := greenlight
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

LDFLAGS := -s -w \
	-X github.com/atlantic-blue/greenlight/internal/version.Version=$(VERSION) \
	-X github.com/atlantic-blue/greenlight/internal/version.GitCommit=$(COMMIT) \
	-X github.com/atlantic-blue/greenlight/internal/version.BuildDate=$(DATE)

.PHONY: build test vet clean

build:
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) .

test:
	go test ./...

vet:
	go vet ./...

clean:
	rm -f $(BINARY)
