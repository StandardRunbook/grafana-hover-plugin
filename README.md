# Hover - Grafana Panel Plugin

[![CI](https://github.com/StandardRunbook/grafana-hover-plugin/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/StandardRunbook/grafana-hover-plugin/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/StandardRunbook/grafana-hover-plugin/branch/main/graph/badge.svg)](https://codecov.io/gh/StandardRunbook/grafana-hover-plugin)
[![Go Report Card](https://goreportcard.com/badge/github.com/StandardRunbook/grafana-hover-plugin)](https://goreportcard.com/report/github.com/StandardRunbook/grafana-hover-plugin)
[![Grafana](https://img.shields.io/badge/Grafana-9.0%2B-orange?logo=grafana)](https://grafana.com)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Latest release](https://img.shields.io/github/v/release/StandardRunbook/grafana-hover-plugin?display_name=tag&sort=semver)](https://github.com/StandardRunbook/grafana-hover-plugin/releases)

A Grafana panel plugin that shows logs from the same time window as the
metric you're hovering on. When the cursor moves over a data point in any
chart on the dashboard, the panel sends the metric name, dashboard, panel
title, and time window to a configured log analysis API and renders the
returned log groups.

The repo also ships a Go backend (`pkg/main.go` → `internal/`) that
implements the API on top of ClickHouse, with Jensen-Shannon divergence
analytics over per-window template histograms.

## How it works

1. Add the Hover panel to a dashboard.
2. Set the API endpoint in the panel options.
3. Hover any chart. The panel reads the metric name, the panel title,
   the dashboard name, and the cursor's chart-time, then POSTs them to
   the API.
4. The API returns log groups; the panel renders them ranked by
   relevance and tagged with a percent-change-from-baseline indicator.

See the [demo bundle](https://github.com/StandardRunbook/grafana-hover-plugin/tree/main/demo)
for an end-to-end script that brings up ClickHouse + Grafana + the
OTel-fed ingest pipeline and lands on a live dashboard with a pre-seeded
incident.

## Installation

### From Grafana Catalog (Recommended - Coming Soon)

```bash
grafana-cli plugins install hover-hover-panel
```

### Manual Installation

1. Download the latest release from [GitHub Releases](https://github.com/StandardRunbook/grafana-hover-plugin/releases)
2. Extract to your Grafana plugins directory:
   ```bash
   # Linux
   unzip hover-hover-panel-*.zip -d /var/lib/grafana/plugins/
   
   # macOS (Homebrew)
   unzip hover-hover-panel-*.zip -d /usr/local/var/lib/grafana/plugins/
   
   # Docker
   # Mount the plugin directory in your docker-compose.yml
   volumes:
     - ./hover-hover-panel:/var/lib/grafana/plugins/hover-hover-panel
   ```
3. Restart Grafana


## Configuration

Set in the panel editor:

| Setting               | Default         | Notes                                                                |
| --------------------- | --------------- | -------------------------------------------------------------------- |
| API Endpoint          | (required)      | URL of the log analysis API. Receives `POST /` with the JSON below. |
| API Key               | empty           | Sent as `Authorization: Bearer <key>` if non-empty.                  |
| Time Window (ms)      | 60000           | Width of the window centered on the hover time.                      |
| Max Logs              | 500             | Cap on log entries rendered.                                         |
| Max Log Length        | 10000           | Per-entry character cap; over-cap entries truncate.                  |
| Log Truncate Length   | 120             | Per-entry threshold for the expand/collapse toggle.                  |

## API Integration

### Request Format

The Hover panel sends POST requests to your API endpoint with this JSON payload:

```json
{
  "org": "1",
  "dashboard": "Production Metrics",
  "panel_title": "CPU Usage",
  "metric_name": "system.cpu.usage",
  "start_time": "2025-01-13T10:30:00.000Z",
  "end_time": "2025-01-13T11:30:00.000Z"
}
```

**Headers:**
- `Content-Type: application/json`
- `Authorization: Bearer <api-key>` (only if API key is configured)

### Expected Response Format

Your API should return log groups with representative logs and change metrics:

```json
{
  "log_groups": [
    {
      "representative_logs": [
        "2025-01-13 11:15:23 ERROR High CPU usage detected",
        "2025-01-13 11:15:24 WARN Thread pool exhausted"
      ],
      "relative_change": 150.5
    }
  ]
}
```

**Fields:**
- `log_groups`: Array of log group objects
- `representative_logs`: Array of log message strings
- `relative_change`: Percentage change from baseline (used for color coding)

### Change Indicators

The panel color-codes log groups based on the `relative_change` value:
- 🔴 Red: > 50% increase (critical)
- 🟠 Orange: 10-50% increase (warning)
- 🟢 Green: < -10% decrease (improvement)
- ⚪ White: -10% to 10% (neutral)

A reference API implementation backed by ClickHouse lives under
`pkg/` + `internal/`. The Go backend reads templates that the
[log_analysis](https://github.com/StandardRunbook/log_analysis) ingest
service has already written and computes JS divergence over per-window
template histograms.

## Requirements

- Grafana 9.0+.
- A log analysis API at the configured endpoint. The bundled Go backend
  uses ClickHouse; supply your own if you'd rather hit Loki / ES / etc.
- The dashboard's "Shared crosshair" setting on, so hover events fire
  on all panels.

## Troubleshooting

**No logs appearing when hovering:**
- Check browser console for API errors
- Verify API endpoint is reachable
- Ensure API returns correct JSON format
- Check CORS settings on your API

**"Waiting for hover data..." message:**
- Move your mouse over a data point in another panel
- Enable shared crosshair in dashboard settings
- Ensure other panels have time-series data

**Plugin not loading:**
- Check Grafana logs for errors
- Verify plugin files are in the correct directory
- For unsigned plugins, add to `allow_loading_unsigned_plugins` in grafana.ini

## Contributing

Contributions are welcome! Please:
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

Apache 2.0 - see [LICENSE](https://github.com/StandardRunbook/grafana-hover-plugin/blob/main/LICENSE) for details.

## Support

- **Issues**: [GitHub Issues](https://github.com/StandardRunbook/grafana-hover-plugin/issues)
- **Discussions**: [GitHub Discussions](https://github.com/StandardRunbook/grafana-hover-plugin/discussions)
- **Author**: Ankil Patel
- **Website**: [ankilp.github.io](https://ankilp.github.io)

## Acknowledgments

Built with [Grafana Plugin Tools](https://grafana.com/developers/plugin-tools/).
