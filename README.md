# Hover - Grafana Panel Plugin

[![Grafana](https://img.shields.io/badge/Grafana-9.0%2B-orange?logo=grafana)](https://grafana.com)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Version](https://img.shields.io/badge/Version-1.0.0-brightgreen)](https://github.com/StandardRunbook/grafana-hover-plugin/releases)

**Automatically correlate metrics with logs when you hover over data points.**

Hover is a Grafana panel plugin that displays related logs when you hover over metrics in other panels on your dashboard. It helps you troubleshoot issues faster by automatically correlating your time-series data with relevant logs from the same time window.


> **🚀 Quick Start**: Add the Hover panel to your dashboard, configure your API endpoint, and start hovering over data points to see correlated logs instantly!

## Features

- **Automatic Log Correlation**: Hover over any data point on your dashboard and instantly see related logs
- **Zero Configuration for Viewing**: Just hover - no clicking required
- **API Integration**: Sends hover context (metric name, time range, panel, dashboard) to your log analysis API
- **Smart Time Windows**: Queries logs within a configurable time window around the hover point
- **Expandable Log Entries**: Long log messages are automatically truncated with click-to-expand
- **Performance Optimized**: Built-in limits for log count and message length
- **Clean UI**: Maximizes space for logs with minimal chrome

## Use Cases

- **Debugging Production Issues**: Quickly correlate metric spikes with error logs
- **Root Cause Analysis**: See what was happening in your logs when metrics changed
- **Incident Response**: Speed up investigation by linking metrics and logs in real-time
- **Performance Troubleshooting**: Find log entries related to slow response times or high resource usage

## How It Works

1. Add the Hover panel to your Grafana dashboard
2. Configure the panel with your log analysis API endpoint
3. Hover over any data point in other panels on the dashboard
4. The Hover panel automatically:
   - Captures the metric name, time, panel title, and dashboard name
   - Queries your API for related logs in that time window
   - Displays the logs grouped by relevance with change indicators

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

The Hover panel is configured through the panel editor UI:

### Required Settings

**API Endpoint** (required)
- The URL of your log analysis API endpoint
- Example: `http://127.0.0.1:3001/query_logs`
- The API receives POST requests with hover context

### Optional Settings

**API Key** (optional)
- API key for authenticating requests to your log analysis API
- If provided, sent as `Authorization: Bearer <api-key>` header
- Leave empty if your API doesn't require authentication

**Time Window (ms)** (default: 3600000 = 1 hour)
- How far back to query logs from the hover time
- Adjust based on your log retention and use case

**Max Logs** (default: 500)
- Maximum number of log entries to display
- Prevents performance issues with large result sets

**Max Log Length** (default: 10000)
- Maximum characters per log entry
- Logs longer than this are truncated

**Log Truncate Length** (default: 120)
- Character threshold for expandable logs
- Logs longer than this show expand/collapse buttons

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

## Example API Implementation

For a complete reference implementation, see the `test-api` directory in the repository.

## Usage Tips

### Enable Shared Crosshair
For the best experience, enable **shared crosshair** in your dashboard settings:
1. Dashboard Settings → General
2. Enable "Shared crosshair"
3. This synchronizes hover events across all panels

### Panel Placement
- Add the Hover panel to a dedicated row at the bottom of your dashboard
- Or use a split view with metrics on the left, Hover panel on the right
- The panel automatically updates as you hover over any other panel

### Performance Tuning
- Adjust **Time Window** based on your needs (smaller = faster queries)
- Set **Max Logs** to limit result size
- Use **Log Truncate Length** to keep the UI clean


## Requirements

- **Grafana**: 9.0.0 or higher
- **API Backend**: You need to provide your own log analysis API
- **CORS**: Your API must allow requests from your Grafana domain

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
