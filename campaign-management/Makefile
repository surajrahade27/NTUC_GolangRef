#include .env if exists
-include .env

# runtime options
COMMIT_HASH = $(shell git rev-parse --short HEAD)
TAG         = $(shell git symbolic-ref -q --short HEAD || git describe --tags --exact-match)
GOPACKAGES  = $(shell go list ./...)
GOFILES		= $(shell find . -type f -name '*.go' -not -path "./vendor/*")
HAS_GOLINT  = $(shell command -v golint)
APP_NAME    = campaign-management
CMD_FOLDER  = "cmd"
API_FOLDER  = "cmd/api/main.go"

BUILD_CMD  ?= api
GOSWAGGER := docker run --rm -it -v $$(pwd):$$(pwd) -w $$(pwd) quay.io/goswagger/swagger:v0.24.0
API_DOCS_PATH = ./swagger.yaml


DC :=
GO       ?= go
LDFLAGS  = -X "main.CommitHash=$(COMMIT_HASH)" -X "main.Tag=$(TAG)"

# HELP
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help

help: ## This help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

# DOCKER TASKS
dc: ## Run targets inside docker-compose
	@echo "Running commands in docker-compose"
	$(eval override DC := docker-compose run --rm app)

.PHONY: tools
tools: ## Install general tools globally (not in the project)
	GO111MODULE=off $(GO) get -u github.com/pressly/goose/cmd/goose
	$(MAKE) test-deps

# Build the container
dc-build: ## Build the container
	docker-compose build

dev-setup:
	git config core.hooksPath .githooks
	chmod +x .githooks/pre-commit
	cp .env.dist .env ## copy env file
	$(MAKE) dc-build ## Setup environment for development

generate-rsa-keys:
	mkdir -p $(KEYPATH)
	openssl genrsa -out $(KEYPATH)/volo.pem 2048
	openssl rsa -in $(KEYPATH)/volo.pem -outform PEM -pubout -out $(KEYPATH)/volo.pub

build: ## Build binaries
	$(DC) $(GO) build -ldflags '$(LDFLAGS)' -o ./bin/api cmd/api/main.go

build-with-arg: ## Build binaries by BUILD_CMD
	$(DC) $(GO) build -ldflags '$(LDFLAGS)' -o ./bin/$(BUILD_CMD) cmd/$(BUILD_CMD)/main.go

build-static: ## Build binaries statically
	$(DC) CGO_ENABLED=0 $(GO) build -ldflags '$(LDFLAGS)' -a -installsuffix cgo -o ./bin/api cmd/api/main.go

tests-verbose: ## Run all tests with verbose flag
	$(DC) $(GO) clean -testcache
	$(DC) $(GO) test -v -race ./app/...

tests: ## Run all tests
	$(DC) $(GO) mod tidy
	$(DC) $(GO) mod vendor
	$(DC) $(GO) clean -testcache
	$(DC) $(GO) test -race ./app/...

tests-e2e: ## Run e2e tests
	$(DC) $(GO) clean -testcache

tests-bench: ## Run benchmark tests
	$(DC) $(GO) test ./bench -bench=.

tests-coverage: ## Run coverage report (excluding mocks)
	$(DC) goverage -race -coverprofile=./cover.out ./app/...
	$(DC) cp cover.out cover.out.tmp && cat cover.out.tmp | grep -v "mocks" > cover.out && rm cover.out.tmp

run: ## Run application
	$(DC) go run $(CMD_FOLDER)

test-deps: ## Install dependencies to run tests
	$(GO) get github.com/stretchr/testify/mock
	GO111MODULE=off $(GO) get -u github.com/schrej/godacov
	GO111MODULE=off $(GO) get -u github.com/haya14busa/goverage

deps: ## Install dependecies for the project
	$(GO) mod download
	$(GO) mod tidy
	$(GO) mod vendor

run-api: ## Run API application
	$(DC) go run $(API_FOLDER)

# Swagger
serve-docs: ## Run Swagger server
	$(GOSWAGGER) serve -F=swagger $(API_DOCS_PATH)

linter: ## Run golangci linter for latest changes
	golangci-lint --new-from-rev=origin/master --config=.golangci.yml run

linter-full: ## Run golangci linter for whole project
	golangci-lint --config=.golangci.yml run

git-init: ## For gosec setup
	git config --global url."https://${BITBUCKET_CREDENTIAL}@bitbucket.org/".insteadOf "https://bitbucket.org/" > /dev/null 2>&1
	
gosec-test: ## Pipeline will be fail when issue found
	go install github.com/securego/gosec/v2/cmd/gosec@latest
	gosec -fmt=sonarqube -out=gosec-report.json ./...

gosec-test-nofail: ## Pipeline not will be fail when issue found, sonar will fail the issue when it exceed the quality gate
	go install github.com/securego/gosec/v2/cmd/gosec@latest
	gosec -fmt=sonarqube -out=gosec-report.json -no-fail ./...