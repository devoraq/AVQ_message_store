# Makefile — Go + golangci-lint (best practice, без .cache в репозитории)

SHELL := /usr/bin/env bash

# --- Tools / versions ---------------------------------------------------------
GO                  ?= go
GOTOOLCHAIN         ?= auto                 # можно зафиксировать, напр. go1.24.9
GOLANGCI_LINT_V     ?= v2.5.0
GOLANGCI_INSTALL_SH ?= https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh
CURL                ?= curl

# --- Paths --------------------------------------------------------------------
# GOBIN извлекаем из окружения Go; если пуст — используем GOPATH/bin
GOBIN := $(shell $(GO) env GOBIN)
ifeq ($(strip $(GOBIN)),)
  GOBIN := $(shell $(GO) env GOPATH)/bin
endif

GOLANGCI_LINT := $(GOBIN)/golangci-lint

# --- Packages / flags ---------------------------------------------------------
PKGS             := ./...
TEST_FLAGS       ?= -count=1
CONC_TEST_FLAGS  ?= -race -count=1 -vet=atomic -shuffle=on
BUILD_DIR        ?= bin
BIN_NAME         ?= app
RACE_BIN         ?= $(BUILD_DIR)/$(BIN_NAME)-race
RUN_MAIN         ?= ./cmd/server/main.go

# для краткости
define _echo
	@printf "\033[1;36m▶ %s\033[0m\n" "$(1)"
endef

.PHONY: help deps install-tools tidy fmt fmt-check vet lint lint-fix lint-verify test test-race cover build build-race run run-race clean clean-caches clean-modcache ci

# --- Help ---------------------------------------------------------------------
help:
	@awk 'BEGIN{FS":.*##"; printf "\nTargets:\n"} /^[a-zA-Z0-9_.-]+:.*##/{printf "  \033[1;32m%-18s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# --- Dependencies -------------------------------------------------------------
deps: ## Download Go modules
	$(call _echo,go mod download)
	@set -euo pipefail; $(GO) mod download

install-tools: ## Install/ensure golangci-lint of pinned version
	$(call _echo,ensure golangci-lint $(GOLANGCI_LINT_V))
	@set -euo pipefail; \
	version_output="$$( { if [[ -x "$(GOLANGCI_LINT)" ]]; then "$(GOLANGCI_LINT)" --version 2>/dev/null || true; fi; } )"; \
	if [[ ! -x "$(GOLANGCI_LINT)" ]] || [[ "$$version_output" != *"$(GOLANGCI_LINT_V)"* ]]; then \
	  install_dir="$(dir $(GOLANGCI_LINT))"; \
	  mkdir -p "$$install_dir"; \
	  $(CURL) -sSfL "$(GOLANGCI_INSTALL_SH)" | sh -s -- -b "$$install_dir" $(GOLANGCI_LINT_V); \
	else \
	  echo "golangci-lint $(GOLANGCI_LINT_V) already installed"; \
	fi

# --- Hygiene ------------------------------------------------------------------
tidy: ## go mod tidy
	$(call _echo,go mod tidy)
	@set -euo pipefail; $(GO) mod tidy

fmt: ## Format code (gofmt)
	$(call _echo,gofmt)
	@$(GO) fmt $(PKGS)

fmt-check: ## Fail if formatting differs
	$(call _echo,gofmt check)
	@set -euo pipefail; \
	out="$$(gofmt -l .)"; \
	if [[ -n "$$out" ]]; then echo "$$out"; echo; echo "Run 'make fmt' to format."; exit 1; fi

vet: ## go vet
	$(call _echo,go vet)
	@$(GO) vet $(PKGS)

# --- Lint ---------------------------------------------------------------------
lint-verify: install-tools ## Verify golangci-lint config before run
	$(call _echo,verify golangci-lint config)
	@set -euo pipefail; \
	echo "Config path: $$($(GOLANGCI_LINT) config path)"; \
	$(GOLANGCI_LINT) config verify

lint: lint-verify deps ## Run golangci-lint
	$(call _echo,golangci-lint run)
	@$(GOLANGCI_LINT) run $(PKGS)

lint-fix: lint-verify deps ## golangci-lint --fix
	$(call _echo,golangci-lint run --fix)
	@$(GOLANGCI_LINT) run --fix $(PKGS)

# --- Tests / coverage ---------------------------------------------------------
test: deps ## Run tests without race detector
	$(call _echo,go test $(TEST_FLAGS))
	@$(GO) test $(TEST_FLAGS) $(PKGS)

test-race: deps ## Run tests with race detector and atomic vet checks
	$(call _echo,go test $(CONC_TEST_FLAGS))
	@$(GO) test $(CONC_TEST_FLAGS) $(PKGS)

cover: deps ## Coverage summary (text) + report file
	$(call _echo,coverage)
	@mkdir -p build
	@$(GO) test -coverprofile=build/coverage.out $(PKGS)
	@$(GO) tool cover -func=build/coverage.out | tail -n 1
	@echo "Coverage file: build/coverage.out"

# --- Build / run --------------------------------------------------------------
build: deps ## Build binary to ./bin/$(BIN_NAME)
	$(call _echo,build $(BUILD_DIR)/$(BIN_NAME))
	@mkdir -p $(BUILD_DIR)
	@$(GO) build -o $(BUILD_DIR)/$(BIN_NAME) ./cmd/server

build-race: deps ## Build binary with race detector
	$(call _echo,build race $(RACE_BIN))
	@mkdir -p $(BUILD_DIR)
	@$(GO) build -race -o $(RACE_BIN) ./cmd/server

run: build ## Run built binary
	$(call _echo,run $(BUILD_DIR)/$(BIN_NAME))
	@./$(BUILD_DIR)/$(BIN_NAME)

run-race: deps ## Run main package with race detector (blocking)
	$(call _echo,go run -race $(RUN_MAIN))
	@$(GO) run -race $(RUN_MAIN)

# --- Clean --------------------------------------------------------------------
clean: ## Remove build artifacts (not caches)
	$(call _echo,clean build artifacts)
	@rm -rf build $(BUILD_DIR)

clean-caches: ## Clean Go & golangci-lint caches (safe)
	$(call _echo,clean caches)
	-@$(GOLANGCI_LINT) cache clean || true
	-@$(GO) clean -cache -testcache || true

clean-modcache: ## Clean module cache (destructive)
	$(call _echo,clean module cache)
	-@$(GO) clean -modcache || true

# --- CI pipeline convenience --------------------------------------------------
ci: tidy fmt-check vet lint test test-race cover ## Run full CI-like pipeline
	$(call _echo,CI pipeline passed)
