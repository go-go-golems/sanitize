.PHONY: all lint lintmax test build goreleaser install

VERSION ?= v0.1.0
GORELEASER_ARGS ?= --skip=sign --snapshot --clean
GORELEASER_TARGET ?= --single-target

lint:
	GOWORK=off golangci-lint run -v

lintmax:
	GOWORK=off golangci-lint run -v --max-same-issues=100

test:
	GOWORK=off go test ./...

build:
	GOWORK=off go build ./...

goreleaser:
	GOWORK=off goreleaser release $(GORELEASER_ARGS) $(GORELEASER_TARGET)

tag-major:
	git tag $(shell svu major)

tag-minor:
	git tag $(shell svu minor)

tag-patch:
	git tag $(shell svu patch)

release:
	git push origin --tags
	GOWORK=off GOPROXY=proxy.golang.org go list -m github.com/go-go-golems/sanitize@$(shell svu current)

SANITIZE_BINARY=$(shell which sanitize 2>/dev/null || echo ./dist/sanitize)
install:
	GOWORK=off go build -o ./dist/sanitize ./cmd/sanitize && \
		cp ./dist/sanitize $(SANITIZE_BINARY)
