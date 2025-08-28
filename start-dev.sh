#!/bin/sh
set -euo pipefail

SERVER_DIR="/portfolio/server"
UI_DIR="/portfolio/ui"

export PNPM_HOME=/portfolio/ui/.pnpm
export PNPM_STORE_PATH=/portfolio/ui/.pnpm-store
export NPM_CONFIG_USERCONFIG=/portfolio/ui/.npmrc
export PATH="$PNPM_HOME:$PATH"
export CI=1

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
  air -c .air.toml 2>&1 | awk '{ cmd="date +%F\\ %T"; cmd | getline ts; close(cmd); printf("[%s] [SERVER] %s\n", ts, $0); fflush(); }'
) &
SERVER_PID=$!

sleep 2

# Start frontend
cd "$UI_DIR" || { log "Failed to enter $UI_DIR"; exit 1; }
log "📦 Installing UI dependencies..."
(
  pnpm install --silent
  if [[ $? -eq 0 ]]; then
    log "[UI] Dependencies successfully installed"
  else
    log "[UI] Error installing dependencies"
    exit 1
  fi
) &

log "🎨 Start the UI development server..."
(
  pnpm run dev -- --host 2>&1 | awk '{ cmd="date +%F\\ %T"; cmd | getline ts; close(cmd); printf("[%s] [UI] %s\n", ts, $0); fflush(); }'
) &
UI_PID=$!

log "✅ Services started:"
log "   - Go Dev server (PID: $SERVER_PID)"
log "   - UI Dev Server (PID: $UI_PID)"
log "📝 Use Ctrl+C to stop all services"

# Wait for either process to exit
wait -n $SERVER_PID $UI_PID
cleanup
