#!/usr/bin/env bash
# One-shot demo orchestrator. Runs end-to-end in ~2 minutes:
#   1. brings up the docker stack (no-op if already running)
#   2. checks the host-side Rust log-ingest-service is reachable
#   3. starts streaming normal traffic in the background (keeps running)
#   4. waits 60s for the analyzer's trailing baseline to accumulate
#   5. injects a 30s incident
#   6. waits 30s for the post-incident phase to flush through OTel
#   7. prints the dashboard URL pointed at the incident window
#
# Designed to be run live during a demo while the presenter narrates —
# each step prints a header so the listener can follow along.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
LOG=/tmp/varied_logs.log

step() { echo; echo "=== $* ==="; }
fail() { echo "ERROR: $*" >&2; exit 1; }

step "1/7  bringing up clickhouse, grafana, otel-collector"
cd "$REPO_ROOT"
# HOST_IP lets the OTel container reach the host-side Rust ingest service.
# `ipconfig getifaddr en0` is the macOS-friendly route; falls back to a sane default.
export HOST_IP=${HOST_IP:-$(ipconfig getifaddr en0 2>/dev/null || echo host.docker.internal)}
echo "  HOST_IP=$HOST_IP"
docker compose up -d
echo -n "  waiting for ClickHouse health..."
until docker exec clickhouse-server clickhouse-client --query "SELECT 1" >/dev/null 2>&1; do
  sleep 1; echo -n "."
done
echo " ok"

step "2/7  checking host-side Rust log-ingest-service is on :4317"
if ! lsof -i :4317 -sTCP:LISTEN >/dev/null 2>&1; then
  fail "nothing is listening on :4317. Start the Rust log-ingest-service from ../log_analysis first."
fi
echo "  ok — OTel collector will stream logs through it to ClickHouse"

step "3/7  starting baseline streamer (normal traffic ~25 lines/sec)"
# Truncate the file so the OTel collector restarts from beginning cleanly
: > "$LOG"
"$SCRIPT_DIR/stream_normal.sh" >/dev/null 2>&1 &
STREAM_PID=$!
echo "  streamer pid=$STREAM_PID — keeps running until you ^C this script or kill it"
trap 'echo "stopping streamer pid=$STREAM_PID"; kill $STREAM_PID 2>/dev/null || true' EXIT

step "4/7  letting the trailing baseline accumulate (60s)"
for s in 60 50 40 30 20 10; do
  echo "  T-${s}s..."
  sleep 10
done

step "5/7  injecting incident (30s of CRITICAL + 5x ERROR + low cpu)"
"$SCRIPT_DIR/incident.sh"

step "6/7  letting post-incident traffic flush through OTel (30s)"
for s in 30 20 10; do
  echo "  T-${s}s..."
  sleep 10
done

step "7/7  ready"
read inc_start inc_end < /tmp/incident_window.txt
from_ms=$(python3 -c "from datetime import datetime,timezone,timedelta;t=datetime.fromisoformat('${inc_start%Z}+00:00');print(int((t-timedelta(minutes=2)).timestamp()*1000))")
to_ms=$(python3 -c "from datetime import datetime,timezone,timedelta;t=datetime.fromisoformat('${inc_end%Z}+00:00');print(int((t+timedelta(minutes=2)).timestamp()*1000))")
echo
echo "  Dashboard:  http://localhost:3000/d/hover-demo/?from=${from_ms}&to=${to_ms}"
echo "  Login:      admin / admin"
echo
echo "  Hover the dip in CPU Usage; the Hover Logs panel ranks"
echo "  'CRITICAL: db pool exhausted...' first with js≈0.13."
echo "  Hover anywhere outside the dip and js drops to <0.001."
echo
echo "  Streamer keeps running so the time range stays fresh —"
echo "  ^C this script when you're done to stop it."

# Block so the streamer keeps running for the live demo
wait
