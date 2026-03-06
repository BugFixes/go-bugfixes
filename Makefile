SHELL := /usr/bin/env bash

GO ?= go
GO_PACKAGES := ./...
TOOLS_BIN := $(CURDIR)/bin
GOLANGCI_LINT := $(TOOLS_BIN)/golangci-lint
GOIMPORTS := $(TOOLS_BIN)/goimports
GOFILES := $(shell find . -type f -name '*.go' -not -path './.cache/*')

export GOCACHE := $(CURDIR)/.cache/go-build
export GOMODCACHE := $(CURDIR)/.cache/go-mod
export GOTMPDIR := $(CURDIR)/.cache/tmp

.PHONY: all
all: check

.PHONY: tools
tools: $(GOLANGCI_LINT) $(GOIMPORTS)

$(GOLANGCI_LINT):
	@mkdir -p "$(TOOLS_BIN)" "$(GOCACHE)" "$(GOMODCACHE)" "$(GOTMPDIR)"
	GOBIN="$(TOOLS_BIN)" $(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

$(GOIMPORTS):
	@mkdir -p "$(TOOLS_BIN)" "$(GOCACHE)" "$(GOMODCACHE)" "$(GOTMPDIR)"
	GOBIN="$(TOOLS_BIN)" $(GO) install golang.org/x/tools/cmd/goimports@latest

.PHONY: tidy
tidy:
	@mkdir -p "$(GOCACHE)" "$(GOMODCACHE)" "$(GOTMPDIR)"
	$(GO) mod tidy

.PHONY: fmt
fmt: tools
	gofmt -w -s $(GOFILES)
	$(GOIMPORTS) -w $(GOFILES)

.PHONY: lint
lint: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run --config configs/golangci.yml

.PHONY: test
test:
	@mkdir -p "$(GOCACHE)" "$(GOMODCACHE)" "$(GOTMPDIR)"
	$(GO) test $(GO_PACKAGES)

.PHONY: test-race
test-race:
	@mkdir -p "$(GOCACHE)" "$(GOMODCACHE)" "$(GOTMPDIR)"
	$(GO) test -race -covermode=atomic -coverprofile=coverage.txt $(GO_PACKAGES)

.PHONY: check
check: lint test-race

.PHONY: clean
clean:
	$(GO) clean ./...
	rm -rf "$(TOOLS_BIN)" "$(CURDIR)/.cache" coverage.txt
