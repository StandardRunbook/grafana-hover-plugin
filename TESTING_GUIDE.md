# Testing Guide - Hover Panel Plugin

This guide provides comprehensive testing procedures for the Hover Panel plugin for Grafana.

## Table of Contents

- [Quick Start Testing](#quick-start-testing)
- [Test Environment Setup](#test-environment-setup)
- [Test Scenarios](#test-scenarios)
- [Backend Testing](#backend-testing)
- [Frontend Testing](#frontend-testing)
- [Integration Testing](#integration-testing)
- [Performance Testing](#performance-testing)
- [Troubleshooting Test Issues](#troubleshooting-test-issues)

---

## Quick Start Testing

### Prerequisites

- Docker and Docker Compose installed
- Plugin built and available in `dist/` directory

### Start Test Environment

```bash
# Build the plugin
pnpm install
pnpm run build

# Start Grafana with the plugin
docker-compose up -d

# View logs
docker-compose logs -f grafana
```

### Access Test Dashboard

1. Open browser: http://localhost:3000
2. Login with default credentials: `admin` / `admin`
3. Navigate to the demo dashboard: **Hover Plugin Demo**

---

## Test Environment Setup

### Option 1: Docker Compose (Recommended)

The repository includes a Docker Compose setup with Grafana pre-configured:

```bash
# Start environment
docker-compose up -d

# Stop environment
docker-compose down

# Restart after changes
pnpm run build && docker-compose restart grafana
```

**Environment Details:**
- Grafana: http://localhost:3000
- Plugin: Pre-loaded as unsigned plugin
- Demo dashboard: Auto-provisioned with test panels
- Backend: Embedded Go binary with mock data

### Option 2: Local Grafana Installation

```bash
# Build plugin
pnpm run build

# Copy to Grafana plugins directory
cp -r dist /var/lib/grafana/plugins/hover-panel

# Configure Grafana to allow unsigned plugins
# Edit grafana.ini:
[plugins]
allow_loading_unsigned_plugins = hover-panel

# Restart Grafana
sudo systemctl restart grafana-server
```

### Option 3: ClickHouse Backend (Production-like Testing)

```bash
# Start ClickHouse
docker run -d -p 9000:9000 -p 8123:8123 --name clickhouse clickhouse/clickhouse-server

# Configure backend
# Edit pkg/config.toml:
[clickhouse]
url = "localhost:9000"
user = "default"
password = ""
database = "default"

# Rebuild and restart
cd pkg
go build -o ../bin/grafana-plugin-api-linux-amd64 ./cmd
cd ..
pnpm run build
docker-compose restart grafana
```

---

## Test Scenarios

### 1. Basic Functionality Tests

#### Test 1.1: Plugin Loads Successfully

**Steps:**
1. Start Grafana
2. Check Grafana logs for plugin registration
3. Navigate to dashboard configuration
4. Add new panel → Select "Hover" panel type

**Expected Results:**
- ✅ Log shows: `Plugin registered pluginId=hover-panel`
- ✅ "Hover" appears in panel type list
- ✅ Panel editor loads without errors

**Verification:**
```bash
docker logs hover-tracker-panel 2>&1 | grep "hover-panel"
```

---

#### Test 1.2: Hover Event Triggers Log Display

**Steps:**
1. Open demo dashboard with Hover panel
2. Hover mouse over a data point in the "CPU Usage" time series panel
3. Observe the Hover Logs panel

**Expected Results:**
- ✅ Hover panel shows "Loading logs..." briefly
- ✅ Mock log data appears with examples:
  - "⚠️  MOCK DATA: ClickHouse database is not connected"
  - "📝 Example anomaly: ERROR: Out of memory on node-3"
- ✅ No errors in browser console

**Browser Console Check:**
```javascript
// Open DevTools → Console
// No errors should appear
```

---

#### Test 1.3: Backend API Responds

**Steps:**
1. Hover over a data point
2. Open browser DevTools → Network tab
3. Filter for "query_logs"
4. Inspect request/response

**Expected Results:**
- ✅ POST request to `/api/plugins/hover-panel/resources/query_logs`
- ✅ Status code: 200 OK
- ✅ Response contains `log_groups` array
- ✅ Response time < 500ms

**Example Response:**
```json
{
  "log_groups": [
    {
      "representative_logs": [
        "⚠️  MOCK DATA: ClickHouse database is not connected",
        "📝 Example anomaly: ERROR: Out of memory on node-3"
      ],
      "relative_change": 2.5,
      "kl_contribution": 0.8,
      "template_id": "mock_error_template"
    }
  ]
}
```

---

### 2. Configuration Tests

#### Test 2.1: Panel Options Are Editable

**Steps:**
1. Edit Hover panel
2. In panel options sidebar, verify all settings are present
3. Modify each setting

**Expected Settings:**
- ✅ Time Window (ms) - default 3600000
- ✅ Max Logs - default 500
- ✅ Max Log Length - default 10000
- ✅ Log Truncate Length - default 120

**Test Changes:**
```
Time Window: 1800000 (30 minutes)
Max Logs: 100
```

**Expected Results:**
- ✅ Settings save without errors
- ✅ Refresh panel to verify persistence

---

#### Test 2.2: Time Window Affects Query

**Steps:**
1. Set Time Window to 3600000 (1 hour)
2. Hover over data point at 12:00
3. Check request payload in Network tab
4. Verify `start_time` and `end_time` span 1 hour

**Expected Payload:**
```json
{
  "start_time": "2025-10-22T11:00:00.000Z",
  "end_time": "2025-10-22T12:00:00.000Z"
}
```

---

### 3. UI/UX Tests

#### Test 3.1: Log Truncation and Expansion

**Steps:**
1. If using mock data, logs are already shown
2. Look for logs with "📝 Example" prefix
3. Check if truncation happens (if log > 120 chars)

**Expected Results:**
- ✅ Long logs show "..." truncation
- ✅ Click to expand shows full text
- ✅ Click again collapses

---

#### Test 3.2: Empty State Display

**Steps:**
1. Stop backend: `docker-compose stop grafana`
2. Or modify backend to return empty log_groups
3. Refresh dashboard

**Expected Results:**
- ✅ Shows "Waiting for hover data..." when not hovering
- ✅ Shows "No logs found" if backend returns empty results

---

#### Test 3.3: Error State Display

**Steps:**
1. Modify backend URL to invalid endpoint (panel settings)
2. Hover over data point

**Expected Results:**
- ✅ Error message displayed in panel
- ✅ Browser console shows error details
- ✅ Panel doesn't crash

---

### 4. Multi-Panel Interaction Tests

#### Test 4.1: Shared Crosshair

**Steps:**
1. Enable shared crosshair in dashboard settings:
   - Dashboard settings → General → Crosshair → Shared crosshair
2. Hover over one time series panel
3. Observe all panels highlight same timestamp

**Expected Results:**
- ✅ Vertical line appears across all panels
- ✅ Hover panel updates when hovering any panel
- ✅ Timestamp synchronization is accurate

---

#### Test 4.2: Multiple Time Series Panels

**Steps:**
1. Create dashboard with 3+ time series panels
2. Hover over each panel individually
3. Verify Hover panel updates for each

**Expected Results:**
- ✅ Panel title changes in request payload
- ✅ Metric name changes based on hovered series
- ✅ No lag or missed events

---

### 5. Data Format Tests

#### Test 5.1: Metric Name Extraction

**Steps:**
1. Hover over time series with different metric names
2. Check Network tab request payload

**Test Cases:**
- Single metric: `cpu_usage`
- Multi-field metric: `system.cpu.usage`
- Label-based metric: `{job="api",instance="server-1"}`

**Expected Results:**
- ✅ Metric name correctly extracted
- ✅ Special characters handled properly

---

#### Test 5.2: Dashboard Context

**Steps:**
1. Rename dashboard
2. Hover over panel
3. Check request payload

**Expected Results:**
- ✅ `dashboard` field contains current dashboard name
- ✅ `org` field contains organization ID
- ✅ `panel_title` matches hovered panel title

---

## Backend Testing

### Test 6: Backend Binary Execution

#### Test 6.1: Binary Permissions

**Steps:**
```bash
ls -la dist/gpx_hover-panel_*
```

**Expected Results:**
- ✅ All binaries have execute permission (`-rwxr-xr-x`)
- ✅ File sizes reasonable (15-20 MB per binary)

---

#### Test 6.2: Backend Startup

**Steps:**
```bash
# Check backend process
docker exec hover-tracker-panel ps aux | grep gpx_hover-panel
```

**Expected Results:**
- ✅ Process is running
- ✅ No zombie processes

**Grafana Logs Check:**
```bash
docker logs hover-tracker-panel 2>&1 | grep "hover-panel"
```

Should show:
```
Plugin registered pluginId=hover-panel
Plugin will return mock data when ClickHouse is unavailable
Successfully started backend plugin process
```

---

#### Test 6.3: Mock Data Fallback

**Steps:**
1. Ensure ClickHouse is NOT running
2. Hover over data point
3. Check response

**Expected Results:**
- ✅ Backend logs: "Plugin will return mock data when ClickHouse is unavailable"
- ✅ Response contains mock log examples
- ✅ No 500 errors

---

#### Test 6.4: ClickHouse Connection (If Available)

**Prerequisites:**
- ClickHouse running
- Tables created (`logs`, `hover_events`)
- Real log data populated

**Steps:**
1. Configure `pkg/config.toml` with ClickHouse credentials
2. Rebuild backend
3. Restart Grafana
4. Hover over data point

**Expected Results:**
- ✅ Backend connects to ClickHouse successfully
- ✅ Real logs returned (not mock data)
- ✅ Logs correlate with hover timestamp
- ✅ No connection errors in Grafana logs

---

## Frontend Testing

### Test 7: React Component Tests

#### Test 7.1: Panel Rendering

**Manual Test:**
```bash
pnpm run test:ci
```

**Expected Results:**
- ✅ All tests pass
- ✅ No console errors
- ✅ Coverage meets thresholds

---

#### Test 7.2: Browser Compatibility

**Test Browsers:**
- Chrome/Chromium (latest)
- Firefox (latest)
- Safari (latest)
- Edge (latest)

**Steps for Each Browser:**
1. Open Grafana dashboard
2. Hover over data points
3. Check console for errors
4. Verify visual rendering

**Expected Results:**
- ✅ Works in all modern browsers
- ✅ No visual glitches
- ✅ Performance acceptable

---

### Test 8: Performance Tests

#### Test 8.1: Large Log Volume

**Setup:**
```typescript
// Simulate backend returning many logs
// In backend: return 500+ log entries
```

**Steps:**
1. Set Max Logs to 1000
2. Hover over data point
3. Measure render time

**Expected Results:**
- ✅ Panel renders in < 2 seconds
- ✅ UI remains responsive
- ✅ Scrolling is smooth

---

#### Test 8.2: Rapid Hover Events

**Steps:**
1. Quickly move mouse across time series panel
2. Observe Hover panel updates
3. Check for memory leaks (DevTools → Performance)

**Expected Results:**
- ✅ Debouncing prevents excessive requests
- ✅ No memory leaks
- ✅ UI doesn't freeze

---

#### Test 8.3: Long Log Messages

**Test Data:**
```
Log message with 10,000+ characters...
```

**Steps:**
1. Configure Max Log Length: 10000
2. Return logs exceeding this length
3. Verify truncation

**Expected Results:**
- ✅ Logs truncated at max length
- ✅ Expand/collapse works
- ✅ No UI overflow

---

## Integration Testing

### Test 9: End-to-End Workflow

#### Test 9.1: Complete User Journey

**Scenario:** DevOps engineer investigating CPU spike

**Steps:**
1. Open dashboard showing CPU metrics
2. Notice spike at 14:30
3. Hover over spike
4. Review logs in Hover panel
5. Identify error logs during spike
6. Expand log entry for details

**Expected Results:**
- ✅ Logs correlate with spike time
- ✅ Error logs clearly visible
- ✅ Full workflow takes < 10 seconds
- ✅ No manual searching required

---

#### Test 9.2: Multi-Metric Investigation

**Steps:**
1. Dashboard with CPU, Memory, Disk I/O panels
2. Hover over CPU spike
3. Review logs
4. Hover over Memory spike
5. Compare logs

**Expected Results:**
- ✅ Each hover returns relevant context
- ✅ Metric names differ in requests
- ✅ Logs help correlate issues

---

### Test 10: Upgrade/Migration Testing

#### Test 10.1: Upgrade from v1.0.3 to v1.0.4

**Steps:**
1. Install v1.0.3
2. Create dashboard with old plugin ID (`hover-hover-panel`)
3. Upgrade to v1.0.4
4. Update dashboard panel type to `hover-panel`

**Expected Results:**
- ✅ Panel still works after upgrade
- ✅ Settings preserved
- ✅ No data loss

---

## Troubleshooting Test Issues

### Common Issues and Solutions

#### Issue: Plugin Not Loading

**Symptoms:**
- Plugin not in panel list
- Error in Grafana logs: "plugin not found"

**Debug Steps:**
```bash
# Check plugin files exist
ls -la dist/

# Check Grafana plugins directory
docker exec hover-tracker-panel ls -la /var/lib/grafana/plugins/hover-panel/

# Check Grafana config
docker exec hover-tracker-panel cat /etc/grafana/grafana.ini | grep unsigned
```

**Solution:**
- Verify plugin ID matches in docker-compose.yaml
- Ensure binaries are executable: `chmod +x dist/gpx_hover-panel_*`
- Check unsigned plugins allowed

---

#### Issue: Backend Not Starting

**Symptoms:**
- "permission denied" error
- "fork/exec" error in logs

**Debug Steps:**
```bash
# Check binary permissions
docker exec hover-tracker-panel ls -la /var/lib/grafana/plugins/hover-panel/gpx_hover-panel_*

# Check backend logs
docker logs hover-tracker-panel 2>&1 | grep -i "error\|permission"
```

**Solution:**
```bash
chmod +x dist/gpx_hover-panel_*
pnpm run build
docker-compose restart grafana
```

---

#### Issue: No Mock Data Shown

**Symptoms:**
- Empty response
- No logs displayed

**Debug Steps:**
```bash
# Check backend error handling
docker logs hover-tracker-panel 2>&1 | grep "mock data"

# Test API directly
curl -X POST http://localhost:3000/api/plugins/hover-panel/resources/query_logs \
  -H "Content-Type: application/json" \
  -d '{"org":"1","dashboard":"test","panel_title":"test","metric_name":"test","start_time":"2025-10-22T10:00:00Z","end_time":"2025-10-22T11:00:00Z"}'
```

**Solution:**
- Verify error matching in `pkg/internal/api/handler.go`
- Ensure "dial tcp" error triggers mock data

---

#### Issue: CORS Errors (If Using External API)

**Symptoms:**
- Browser console: "CORS policy" error
- Network tab shows preflight failure

**Solution:**
```go
// In your API backend, add CORS headers
w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
```

---

## Test Checklist

Use this checklist before releases:

### Pre-Release Testing

- [ ] Plugin loads in Grafana
- [ ] Backend binary starts successfully
- [ ] Mock data displays when ClickHouse unavailable
- [ ] Hover events trigger log queries
- [ ] Network requests complete successfully (200 OK)
- [ ] Panel settings are editable and persist
- [ ] Log truncation and expansion works
- [ ] Shared crosshair synchronizes panels
- [ ] Multiple panel types supported (time series, graph, etc.)
- [ ] Browser console has no errors
- [ ] All platform binaries executable
- [ ] Plugin signature valid
- [ ] Demo dashboard loads correctly

### Performance Checklist

- [ ] Response time < 500ms for typical queries
- [ ] UI responsive with 500+ logs
- [ ] No memory leaks during extended use
- [ ] Rapid hover events handled gracefully

### Browser Compatibility

- [ ] Chrome/Chromium
- [ ] Firefox
- [ ] Safari
- [ ] Edge

### Documentation

- [ ] README accurate
- [ ] API specification up to date
- [ ] Configuration examples valid
- [ ] Troubleshooting guide helpful

---

## Automated Testing

### Unit Tests

```bash
pnpm run test
```

### E2E Tests (Future)

```bash
pnpm run e2e
```

### Linting

```bash
pnpm run lint
```

### Type Checking

```bash
pnpm run typecheck
```

---

## Reporting Issues

When reporting test failures, include:

1. **Environment Details:**
   - Grafana version
   - Plugin version
   - OS/Platform
   - Browser (if frontend issue)

2. **Reproduction Steps:**
   - Step-by-step instructions
   - Sample data/configuration

3. **Expected vs Actual:**
   - What should happen
   - What actually happened

4. **Logs/Screenshots:**
   - Browser console errors
   - Grafana logs
   - Network tab screenshots

5. **Configuration:**
   - Panel settings
   - Backend config (sanitized)

**Submit to:** https://github.com/StandardRunbook/grafana-hover-plugin/issues

---

## Test Data Generation

### Mock ClickHouse Data

```sql
-- Create test tables
CREATE TABLE logs (
  timestamp DateTime,
  level String,
  message String,
  service String
) ENGINE = MergeTree()
ORDER BY timestamp;

-- Insert test data
INSERT INTO logs VALUES
  ('2025-10-22 10:30:00', 'ERROR', 'Out of memory', 'api-server'),
  ('2025-10-22 10:30:01', 'WARN', 'High CPU usage', 'api-server'),
  ('2025-10-22 10:30:05', 'INFO', 'Request processed', 'api-server');
```

---

## Performance Benchmarks

### Target Metrics

| Metric | Target | Acceptable | Critical |
|--------|--------|------------|----------|
| API Response Time | < 200ms | < 500ms | > 1s |
| Panel Render Time | < 500ms | < 1s | > 2s |
| Memory Usage | < 50MB | < 100MB | > 200MB |
| Log Display (500 logs) | < 1s | < 2s | > 3s |

---

## Continuous Testing

### GitHub Actions (Future)

```yaml
# .github/workflows/test.yml
name: Test Plugin
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - run: pnpm install
      - run: pnpm run test:ci
      - run: pnpm run lint
      - run: pnpm run typecheck
```

---

## Contact

For testing questions or issues:
- GitHub Issues: https://github.com/StandardRunbook/grafana-hover-plugin/issues
- Author: Ankil Patel
