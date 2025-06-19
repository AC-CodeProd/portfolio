#!/bin/sh
set -euo pipefail

SERVER_DIR="/portfolio/server"
UI_DIR="/portfolio/ui"

log() {
  echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1"
}

command -v air >/dev/null 2>&1 || { log >&2 "air is not installed. Aborting."; exit 1; }
command -v pnpm >/dev/null 2>&1 || { log >&2 "pnpm is not installed. Aborting."; exit 1; }


cleanup() {
  log "Script stopping, cleaning up..."
  pkill -P $$ || true
  exit 0
}
trap cleanup INT TERM

# Start backend server
cd "$SERVER_DIR" || { log "Failed to enter $SERVER_DIR"; exit 1; }
log "🔧 Starting the Go development server..."
(
  cd "$SERVER_DIR"
  air -c .air.toml 2>&1 | while IFS= read -r line; do
    echo "[SERVER] $line"
  done
) &
SERVER_PID=$!

sleep 2

# Start frontend
cd "$UI_DIR" || { log "Failed to enter $UI_DIR"; exit 1; }
log "📦 Installing UI dependencies..."
(
  pnpm install --silent
  if [[ $? -eq 0 ]]; then
    echo "[UI] Dependencies successfully installed"
  else
    echo "[UI] Error installing dependencies"
    exit 1
  fi
) &

log "🎨 Start the UI development server..."
(
  pnpm run dev 2>&1 | while IFS= read -r line; do
    echo "[UI] $line"
  done
) &
UI_PID=$!

log "✅ Services started:"
log "   - Go Dev server (PID: $SERVER_PID)"
log "   - UI Dev Server (PID: $UI_PID)"
log "📝 Use Ctrl+C to stop all services"

# Wait for either process to exit
wait -n $SERVER_PID $UI_PID
cleanup
