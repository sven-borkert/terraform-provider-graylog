GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
BINARY := terraform-provider-graylog
BIN_DIR := bin
PROVIDER_PATH := ./cmd/terraform-provider-graylog
MIRROR_DIR := $(HOME)/.terraform.d/plugins/registry.terraform.io/sven-borkert/graylog/0.0.0/$(GOOS)_$(GOARCH)
TF_DEV_PROVIDER_DIR ?= $(HOME)/.terraform-providers
TF_CLI_CONFIG_PATH ?= $(HOME)/.terraformrc
TF_LOCAL_CONFIG_SCRIPT := ./bin/configure-terraform-local-provider
TF_LOCAL_PROVIDER_BINARY := $(TF_DEV_PROVIDER_DIR)/$(BINARY)
VENV := .venv

.PHONY: build dev-install local-provider-build local-provider-config local-provider-setup clean venv venv-install

build:
	go build -o $(BIN_DIR)/$(BINARY) $(PROVIDER_PATH)

dev-install: build
	mkdir -p $(MIRROR_DIR)
	cp $(BIN_DIR)/$(BINARY) $(MIRROR_DIR)/$(BINARY)_v0.0.0
	@echo "Installed to $(MIRROR_DIR)"

local-provider-build:
	mkdir -p $(TF_DEV_PROVIDER_DIR)
	go build -o $(TF_LOCAL_PROVIDER_BINARY) $(PROVIDER_PATH)
	@echo "Installed local provider to $(TF_LOCAL_PROVIDER_BINARY)"

local-provider-config:
	$(TF_LOCAL_CONFIG_SCRIPT) "$(TF_DEV_PROVIDER_DIR)" "$(TF_CLI_CONFIG_PATH)"
	@echo "Configured Terraform CLI config at $(TF_CLI_CONFIG_PATH)"

local-provider-setup: local-provider-build local-provider-config

clean:
	rm -rf $(BIN_DIR)/*

venv:
	python3 -m venv $(VENV)
	$(VENV)/bin/pip install --upgrade pip

venv-install: venv
	$(VENV)/bin/pip install requests
	@echo "Venv ready at $(VENV)"

.PHONY: fmt lint test acc-test

fmt:
	@if command -v gofumpt >/dev/null 2>&1; then \
	  gofumpt -l -s -w $$(git ls-files '*.go'); \
	else \
	  gofmt -w $$(git ls-files '*.go'); \
	fi

lint:
	golangci-lint run

test:
	go test -race -covermode=atomic ./...

acc-test:
	TF_ACC=1 go test -v -race -covermode=atomic ./graylog/resource/... ./graylog/datasource/...
