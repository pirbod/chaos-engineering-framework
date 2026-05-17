#!/usr/bin/env bash
set -euo pipefail

BACKEND_PORT="${BACKEND_PORT:-8080}"
FRONTEND_PORT="${FRONTEND_PORT:-5173}"

./scripts/seed-demo-data.sh

if ! command -v go >/dev/null 2>&1; then
  echo "Go is required for make dev. Use docker-compose up --build if Go is not installed."
  exit 1
fi

if [ ! -d "frontend/node_modules" ]; then
  echo "Installing frontend dependencies..."
  (cd frontend && npm install)
fi

echo "Backend:  http://localhost:${BACKEND_PORT}"
echo "Frontend: http://localhost:${FRONTEND_PORT}"

trap 'jobs -p | xargs -r kill' EXIT

(cd backend && PORT="${BACKEND_PORT}" VERSION="local-dev" go run ./cmd/server) &
(cd frontend && VITE_API_BASE_URL="http://localhost:${BACKEND_PORT}" npm run dev -- --port "${FRONTEND_PORT}") &

wait
