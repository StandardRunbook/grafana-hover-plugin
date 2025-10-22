# Grafana Provisioning

This directory contains Grafana provisioning configuration for automatically setting up the Hover plugin demo environment.

## What's Included

### Dashboards (`dashboards/`)
- **dashboards.yaml**: Dashboard provider configuration
- **hover-demo.json**: Demo dashboard with:
  - 2 time-series panels (CPU and Memory usage with test data)
  - 1 Hover panel configured to display logs on hover
  - Shared crosshair enabled for synchronized hover events

### Data Sources (`datasources/`)
- **testdata.yaml**: TestData data source for generating sample metrics

## How to Use

### With Docker Compose

The provisioning directory is automatically mounted in `docker-compose.test.yml`:

```bash
docker-compose -f docker-compose.test.yml up
```

Then:
1. Open Grafana at http://localhost:3000
2. Login with `admin` / `admin`
3. The "Hover Plugin Demo" dashboard will be available
4. Hover over the CPU or Memory panels to see logs in the Hover panel

### Manual Setup

Copy the provisioning files to your Grafana instance:

```bash
# Copy to Grafana provisioning directory
cp -r provisioning/* /etc/grafana/provisioning/

# Or for Docker
docker cp provisioning/. grafana:/etc/grafana/provisioning/

# Restart Grafana
docker restart grafana
```

## Demo Dashboard

The demo dashboard includes:
- **CPU Usage Panel**: Random walk metric simulating CPU usage
- **Memory Usage Panel**: Random walk metric simulating memory usage
- **Hover Logs Panel**: Configured to call `http://localhost:3001/query_logs`

### Prerequisites for Full Demo

The Hover panel requires the mock API server:

```bash
cd test-api
npm install
node server.js
```

The mock API provides sample log responses when you hover over metrics.

## Customization

### Change API Endpoint

Edit `dashboards/hover-demo.json` and update the `apiEndpoint` option:

```json
"options": {
  "apiEndpoint": "http://your-api-server:port/endpoint",
  "timeWindowMs": 3600000,
  "maxLogs": 500
}
```

### Add Your Own Dashboards

1. Create a JSON file in `dashboards/`
2. Include the Hover panel with your configuration
3. Restart Grafana or wait for the provisioning interval

### Add Data Sources

Create YAML files in `datasources/` following this format:

```yaml
apiVersion: 1
datasources:
  - name: MyDataSource
    type: prometheus  # or other type
    url: http://prometheus:9090
    access: proxy
```

## Testing

To verify provisioning is working:

```bash
# Check Grafana logs
docker logs hover-plugin-grafana

# Look for provisioning messages:
# "Inserting datasource from configuration"
# "Provisioned dashboards from..."
```

## Learn More

- [Grafana Provisioning Docs](https://grafana.com/docs/grafana/latest/administration/provisioning/)
- [Dashboard Provisioning](https://grafana.com/docs/grafana/latest/administration/provisioning/#dashboards)
- [Data Source Provisioning](https://grafana.com/docs/grafana/latest/administration/provisioning/#data-sources)
