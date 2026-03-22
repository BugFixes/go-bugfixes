set shell := ["bash", "-uc"]

go := env_var_or_default("GO", "go")
root := justfile_directory()
tools_bin := root / "bin"
golangci_lint := tools_bin / "golangci-lint"
goimports := tools_bin / "goimports"
go_packages := "./..."
go_cache := root / ".cache/go-build"
go_mod_cache := root / ".cache/go-mod"
go_tmp := root / ".cache/tmp"

export GOCACHE := go_cache
export GOMODCACHE := go_mod_cache
export GOTMPDIR := go_tmp

default: check

_prepare:
    mkdir -p "{{ tools_bin }}" "{{ go_cache }}" "{{ go_mod_cache }}" "{{ go_tmp }}"

tools: _prepare
    GOBIN="{{ tools_bin }}" {{ go }} install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    GOBIN="{{ tools_bin }}" {{ go }} install golang.org/x/tools/cmd/goimports@latest

tidy: _prepare
    {{ go }} mod tidy

fmt: tools
    find . -type f -name '*.go' -not -path './.cache/*' -print0 | xargs -0 gofmt -w -s
    find . -type f -name '*.go' -not -path './.cache/*' -print0 | xargs -0 "{{ goimports }}" -w

lint: tools
    "{{ golangci_lint }}" run --config configs/golangci.yml

test: _prepare
    {{ go }} test {{ go_packages }}

test-race: _prepare
    {{ go }} test -race -covermode=atomic -coverprofile=coverage.txt {{ go_packages }}

check: lint test-race

clean:
    {{ go }} clean ./...
    rm -rf "{{ tools_bin }}" "{{ root }}/.cache" coverage.txt
