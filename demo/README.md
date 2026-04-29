# Demo

End-to-end demo of the hover plugin: hover any time on a metric chart →
see the actual log lines responsible for that moment, ranked by
JS-divergence against a trailing 2-hour baseline.

## Prerequisites

- Docker Desktop (for the ClickHouse / Grafana / OTel collector stack)
- Python 3 (for the timestamp math in `run.sh`)
- The Rust `log-ingest-service` from `../log_analysis` running on
  `host:4317`. The demo script verifies this and exits early if it's
  missing.

## One-shot

```bash
./demo/run.sh
```

Total time ~2 minutes. The script narrates each phase so you can talk
over it. When it finishes it prints a dashboard URL pointing directly
at the engineered incident.

## What it does

1. `docker compose up -d` — starts ClickHouse, Grafana with the plugin
   mounted from `./dist`, and the OTel collector pointing at the host's
   Rust ingest service.
2. Truncates `/tmp/varied_logs.log` and starts streaming mixed normal
   traffic (cpu/memory/User/disk/ERROR, ~25 lines/sec). Keeps running
   in the background until you ^C the script.
3. Sleeps 60 s so the analyzer's trailing baseline has data to compare
   against.
4. Calls `incident.sh` — 30 s of distribution-shifted traffic: cpu
   drops to 5–15 %, ERROR rate 5×, and a brand-new
   `CRITICAL: db pool exhausted…` template appears that didn't exist
   in baseline.
5. Sleeps 30 s for the post-incident phase to flush.
6. Prints a deep-link to the dashboard with the time range pre-set
   to the incident window ± 2 minutes.

## What you should see

- **CPU Usage panel:** sharp dip from ~70 % to ~10 % in the middle.
- **Error Rate panel:** spike co-located with the dip.
- **Hover Logs panel** when you hover the dip: top entry is
  `CRITICAL: db pool exhausted…` with `js≈0.13` and `Δ≈+1200%`.
- **Hover Logs panel** when you hover the flat areas either side: top
  entries are baseline templates with `js < 0.001` and `Δ ≈ 0%`.

## Just-the-incident (skip the cold-start)

If the stack is already up and a baseline streamer is already running:

```bash
./demo/incident.sh
# then read /tmp/incident_window.txt for the absolute time window
```

## Cleanup

```bash
docker compose down -v          # tears down ClickHouse data too
pkill -f demo/stream_normal.sh  # if a streamer is still running
```
