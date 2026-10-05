#!/usr/bin/env bash
set -euo pipefail

# Colors (optional)
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log() { echo -e "${GREEN}[run]${NC} $1"; }
warn() { echo -e "${YELLOW}[run]${NC} $1"; }

# Run tests first
log "Running tests..."
go test ./... -race -count=1

# Run the server.
PORT="${PORT:-8080}"
DB_PATH="${DB_PATH:-pismo.db}"

log "Starting API on :${PORT} (db=${DB_PATH})"
PORT="$PORT" DB_PATH="$DB_PATH" go run ./cmd/api