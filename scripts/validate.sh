#!/usr/bin/env bash
set -euo pipefail

run_or_skip() {
  local tool="$1"
  shift
  if command -v "${tool}" >/dev/null 2>&1; then
    "$@"
  else
    echo "${tool} is not installed; skipping: $*"
  fi
}

run_or_skip go bash -lc "cd backend && go test ./..."
run_or_skip gofmt bash -lc "test -z \"\$(gofmt -l backend)\""

if [ ! -d "frontend/node_modules" ]; then
  (cd frontend && npm install)
fi
(cd frontend && npm run lint && npm test && npm run build)

run_or_skip terraform terraform -chdir=terraform fmt -check -recursive
run_or_skip helm helm lint helm/platform-chaos-hub
run_or_skip helm bash -lc "helm template platform-chaos-hub helm/platform-chaos-hub >/tmp/platform-chaos-hub.yaml"

echo "Validation complete."
