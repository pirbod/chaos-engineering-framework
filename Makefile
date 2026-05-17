SHELL := /usr/bin/env bash

BACKEND_DIR := backend
FRONTEND_DIR := frontend

.PHONY: setup dev test lint build docker demo validate

setup:
	@if command -v go >/dev/null 2>&1; then cd $(BACKEND_DIR) && go mod download; else echo "Go is not installed; backend setup will run where Go is available."; fi
	cd $(FRONTEND_DIR) && npm install

dev:
	./scripts/demo.sh

test:
	@if command -v go >/dev/null 2>&1; then cd $(BACKEND_DIR) && go test ./...; else echo "Go is not installed; skipping backend tests in this environment."; fi
	cd $(FRONTEND_DIR) && npm test

lint:
	@if command -v gofmt >/dev/null 2>&1; then test -z "$$(gofmt -l $(BACKEND_DIR))"; else echo "gofmt is not installed; skipping Go formatting check."; fi
	cd $(FRONTEND_DIR) && npm run lint
	@if command -v terraform >/dev/null 2>&1; then terraform -chdir=terraform fmt -check -recursive; else echo "terraform is not installed; skipping Terraform fmt."; fi
	@if command -v helm >/dev/null 2>&1; then helm lint helm/platform-chaos-hub; else echo "helm is not installed; skipping Helm lint."; fi

build:
	@if command -v go >/dev/null 2>&1; then cd $(BACKEND_DIR) && go build -o bin/platform-chaos-hub ./cmd/server; else echo "Go is not installed; skipping backend build."; fi
	cd $(FRONTEND_DIR) && npm run build

docker:
	docker compose build

demo:
	./scripts/demo.sh

validate:
	./scripts/validate.sh
