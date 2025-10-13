# Grafana Hover Tracker Panel

Track and visualize hover actions across your Grafana dashboard.

## Features

- Real-time tracking of hover events across all panels in a dashboard
- Visual timeline of hover activity
- Configurable history size
- Cross-panel communication using BroadcastChannel API
- Shows panel name, element type, and coordinates
- Optional timestamp display

## Installation

1. Clone this repository into your Grafana plugins directory:
```bash
cd /var/lib/grafana/plugins
git clone <repository-url> hover-tracker-panel
```

2. Install dependencies and build:
```bash
cd hover-tracker-panel
pnpm install
pnpm run build
```

3. Restart Grafana

## Usage

1. Add the "Hover Tracker Panel" to your dashboard
2. Configure the panel options:
   - **History Size**: Number of hover events to keep (default: 100)
   - **Show Timestamp**: Display timestamp for each event
   - **Track Own Panel**: Include hover events from the tracker panel itself
   - **Event Channel**: Custom channel name for multi-dashboard tracking

## Development

### Option 1: Docker Compose (Recommended)

```bash
# Install dependencies and build
pnpm install
pnpm run dev

# In another terminal, run Grafana
pnpm run server
```

Access Grafana at http://localhost:3000 (admin/admin)

### Option 2: Local Grafana Installation

1. Build the plugin:
```bash
pnpm install
pnpm run build
```

2. Copy or symlink to your Grafana plugins directory:
```bash
# macOS
ln -s $(pwd)/dist /usr/local/var/lib/grafana/plugins/hover-tracker-panel

# Linux
sudo ln -s $(pwd)/dist /var/lib/grafana/plugins/hover-tracker-panel
```

3. Add to Grafana config:
```ini
[plugins]
allow_loading_unsigned_plugins = hover-tracker-panel
```

4. Restart Grafana

### Option 3: Development Mode

```bash
# Terminal 1: Run webpack in watch mode
pnpm run dev

# Terminal 2: Run Grafana with live reload
docker run -p 3000:3000 \
  -e "GF_PLUGINS_ALLOW_LOADING_UNSIGNED_PLUGINS=hover-tracker-panel" \
  -v $(pwd)/dist:/var/lib/grafana/plugins/hover-tracker-panel \
  grafana/grafana:latest
```

## How it Works

The plugin uses the BroadcastChannel API to communicate hover events between panels. It attaches a global event listener to the dashboard that captures hover events on any panel and broadcasts them to all listening panels.

Each hover event includes:
- Panel ID and title
- Mouse coordinates
- Element type (chart, table, header, etc.)
- Timestamp

## Configuration Options

- **History Size**: Controls how many events are stored in memory
- **Show Timestamp**: Toggle timestamp visibility
- **Track Own Panel**: Whether to include hover events from the tracker panel
- **Event Channel**: Allows multiple tracker panels with different channels

## Important

- Add shared crosshairs from dashboard settings
