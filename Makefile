-include .env

TARGETS := playground

#VERSION := $(shell git describe --tags)
VERSION := 0.0.1
BUILD := $(shell git rev-parse --short HEAD)
DATE := $(shell date +%Y-%m-%dT%TZ%z)
PROJECTNAME := $(shell basename "$(PWD)")

# Go
GOPROXY := https://goproxy.cn
GOBIN := $(shell pwd)/bin
GOFILES := $(wildcard *.go)
GOMODULE := github.com/master-g/playground

GOLANGCILINT := $(GOBIN)/golangci-lint

# output
BIN := $(shell pwd)/bin

# Redirect error output
# STDERR := /tmp/.$(PROJECTNAME)-stderr.txt

# Make is verbose in Linux. Make it silent.
# MAKEFLAGS += --silent


## mod: Recreate go mod files
.PHONY: mod
mod:
	@echo "  >  Recreating go.mod..."
	@rm go.mod
	@rm go.sum
	@go mod init $(GOMODULE)


## vendor: Module cleanup and vendor
.PHONY: vendor
vendor:
	@echo "  >  Module tidy and vendor..."
	@GOPROXY=$(GOPROXY) go mod tidy
	@GOPROXY=$(GOPROXY) go mod download


## lint: Lint go source files
$(GOLANGCILINT):
	@GOBIN=$(GOBIN) wget -O - -q https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s v1.21.0


.PHONY: lint
lint: $(GOLANGCILINT)
	@echo "  >  Linting..."
	@$(GOBIN)/golangci-lint run ./...


## fmt: Formats go source files
.PHONY: fmt
fmt:
	@echo "  >  Formating..."
	@find . -type f -name '*.go' -not -path './vendor/*' -not -path './pb/*' -not -path './.idea/*' -print0 | xargs -0 goimports -w


## build: Build all executables
.PHONY: build
build: $(TARGETS)


## clean: Cleaning build cache
.PHONY: clean
clean:
	@echo "  >  Cleaning build cache..."
	@-for target in $(TARGETS); do rm -f $(BIN)/$$target; done;


$(TARGETS):
	@echo "  >  Building $@..."
	@-go build -ldflags \
	"-X $(GOMODULE)/cmd/$@/cmd.Version=$(VERSION) \
	-X $(GOMODULE)/cmd/$@/cmd.BuildDate=$(DATE) \
	-X $(GOMODULE)/cmd/$@/cmd.CommitHash=$(BUILD)" \
	-o $(BIN)/$@ ./cmd/$@

## release: Build executable files for macOS, Windows
.PHONY: release
release:
	@echo "  >  Releasing..."
	@GOBIN=$(GOBIN) \
	gox -ldflags \
	"-X $(GOMODULE)/pkg/version.Version=$(VERSION) \
	-X $(GOMODULE)/pkg/version.BuildDate=$(DATE) \
	-X $(GOMODULE)/pkg/version.CommitHash=$(BUILD)" \
	-osarch="darwin/amd64" \
	-osarch="windows/amd64" \
	-osarch="linux/amd64" \
	-output="release/{{.OS}}_{{.Arch}}/{{.Dir}}" ./cmd/oxytocin


.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
