# Grafana Hover Tracker Panel - API Specification

## Overview
This document specifies the API requirements for the Grafana Hover Tracker Panel plugin. The plugin sends hover event data to a configured API endpoint and expects log analysis results in return.

## API Endpoint Specification

### Request Details
- **Method**: `POST`
- **Content-Type**: `application/json`
- **Authentication**: Optional Bearer token via `Authorization` header
- **Default Endpoint**: `http://127.0.0.1:3001/query_logs`

### Request Payload Structure

```json
{
  "org": "1",
  "dashboard": "Production Dashboard",
  "panel_title": "CPU Usage",
  "metric_name": "cpu_usage_percent",
  "start_time": "2024-01-15T10:30:00.000Z",
  "end_time": "2024-01-15T11:30:00.000Z"
}
```

#### Field Descriptions

| Field | Type | Description | Example |
|-------|------|-------------|---------|
| `org` | string | Grafana organization ID | "1" |
| `dashboard` | string | Name of the dashboard containing the hovered panel | "Production Dashboard" |
| `panel_title` | string | Title of the panel that was hovered | "CPU Usage" |
| `metric_name` | string | Name of the metric/series that was hovered | "cpu_usage_percent" |
| `start_time` | string | ISO 8601 timestamp for the start of the time window | "2024-01-15T10:30:00.000Z" |
| `end_time` | string | ISO 8601 timestamp for the end of the time window | "2024-01-15T11:30:00.000Z" |

### Expected Response Format

```json
{
  "log_groups": [
    {
      "representative_logs": [
        "2024-01-15T10:45:23Z [ERROR] Database connection timeout",
        "2024-01-15T10:45:24Z [WARN] Retrying connection attempt 1",
        "2024-01-15T10:45:25Z [INFO] Connection established successfully"
      ],
      "relative_change": 15.5
    },
    {
      "representative_logs": [
        "2024-01-15T10:50:12Z [ERROR] Memory allocation failed",
        "2024-01-15T10:50:13Z [WARN] Falling back to alternative memory pool"
      ],
      "relative_change": -8.2
    }
  ]
}
```

#### Response Field Descriptions

| Field | Type | Description |
|-------|------|-------------|
| `log_groups` | array | Array of log groups found during the time window |
| `representative_logs` | array of strings | Array of log entries representing each group |
| `relative_change` | number | Percentage change from baseline (positive or negative) |

## Plugin Configuration Options

The plugin can be configured with the following parameters:

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `apiEndpoint` | string | `http://127.0.0.1:3001/query_logs` | URL endpoint to send hover data |
| `apiKey` | string | `""` | Optional API key for authentication |
| `timeWindowMs` | number | `3600000` | Time window for log queries in milliseconds (1 hour) |
| `maxLogs` | number | `500` | Maximum number of log entries to display |
| `maxLogLength` | number | `10000` | Maximum length of individual log entries |
| `logTruncateLength` | number | `120` | Character threshold for expandable logs |

## Error Handling

### Expected Behavior
- **Success**: HTTP 200 status with JSON response as specified above
- **Error**: Any non-200 status code will be treated as an error
- **Timeout**: Plugin will handle network timeouts gracefully
- **Invalid Response**: Plugin will log warnings for unexpected response formats

### Error Response Example
```json
{
  "error": "Invalid time range",
  "message": "Start time must be before end time",
  "code": 400
}
```

## Implementation Example

### Node.js/Express Example
```javascript
const express = require('express');
const app = express();

app.use(express.json());

app.post('/query_logs', async (req, res) => {
  try {
    const { org, dashboard, panel_title, metric_name, start_time, end_time } = req.body;
    
    // Validate required fields
    if (!org || !dashboard || !panel_title || !metric_name || !start_time || !end_time) {
      return res.status(400).json({
        error: 'Missing required fields',
        message: 'All fields are required'
      });
    }
    
    // Your log analysis logic here
    const logGroups = await analyzeLogs({
      org,
      dashboard, 
      panelTitle: panel_title,
      metricName: metric_name,
      startTime: new Date(start_time),
      endTime: new Date(end_time)
    });
    
    res.json({
      log_groups: logGroups
    });
  } catch (error) {
    console.error('Error processing log query:', error);
    res.status(500).json({
      error: 'Internal server error',
      message: 'Failed to process log query'
    });
  }
});

async function analyzeLogs({ org, dashboard, panelTitle, metricName, startTime, endTime }) {
  // Implement your log analysis logic here
  // This should return an array of log groups with representative_logs and relative_change
  
  return [
    {
      representative_logs: [
        "Sample log entry 1",
        "Sample log entry 2"
      ],
      relative_change: 15.5
    }
  ];
}
```

### Python/Flask Example
```python
from flask import Flask, request, jsonify
from datetime import datetime
import logging

app = Flask(__name__)

@app.route('/query_logs', methods=['POST'])
def query_logs():
    try:
        data = request.get_json()
        
        # Validate required fields
        required_fields = ['org', 'dashboard', 'panel_title', 'metric_name', 'start_time', 'end_time']
        for field in required_fields:
            if field not in data:
                return jsonify({
                    'error': 'Missing required fields',
                    'message': f'Field {field} is required'
                }), 400
        
        # Your log analysis logic here
        log_groups = analyze_logs(
            org=data['org'],
            dashboard=data['dashboard'],
            panel_title=data['panel_title'],
            metric_name=data['metric_name'],
            start_time=datetime.fromisoformat(data['start_time'].replace('Z', '+00:00')),
            end_time=datetime.fromisoformat(data['end_time'].replace('Z', '+00:00'))
        )
        
        return jsonify({
            'log_groups': log_groups
        })
        
    except Exception as e:
        logging.error(f'Error processing log query: {e}')
        return jsonify({
            'error': 'Internal server error',
            'message': 'Failed to process log query'
        }), 500

def analyze_logs(org, dashboard, panel_title, metric_name, start_time, end_time):
    # Implement your log analysis logic here
    # This should return a list of dictionaries with 'representative_logs' and 'relative_change'
    
    return [
        {
            'representative_logs': [
                'Sample log entry 1',
                'Sample log entry 2'
            ],
            'relative_change': 15.5
        }
    ]
```

## Testing the API

### Test Request
```bash
curl -X POST http://localhost:3001/query_logs \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-api-key" \
  -d '{
    "org": "1",
    "dashboard": "Test Dashboard",
    "panel_title": "Test Panel",
    "metric_name": "test_metric",
    "start_time": "2024-01-15T10:30:00.000Z",
    "end_time": "2024-01-15T11:30:00.000Z"
  }'
```

### Expected Test Response
```json
{
  "log_groups": [
    {
      "representative_logs": [
        "2024-01-15T10:45:23Z [INFO] Test log entry 1",
        "2024-01-15T10:45:24Z [WARN] Test log entry 2"
      ],
      "relative_change": 0.0
    }
  ]
}
```

## Notes for Implementation

1. **Time Format**: All timestamps are in ISO 8601 format with 'Z' suffix (UTC)
2. **Log Analysis**: The `relative_change` field should represent the percentage change from a baseline
3. **Performance**: Consider implementing caching for frequently requested time windows
4. **Rate Limiting**: Consider implementing rate limiting to prevent abuse
5. **Authentication**: The API key is optional but recommended for production use
6. **Logging**: Implement proper logging for debugging and monitoring

## Contact

For questions about this API specification, please contact the plugin development team or refer to the plugin documentation.
