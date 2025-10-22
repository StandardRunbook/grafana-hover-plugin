# Changelog

All notable changes to the Grafana Hover Tracker Panel plugin will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.1] - 2025-10-22

### Added
- **Backend binaries included**: Plugin now includes Rust-based backend API server for all platforms
  - Linux AMD64 and ARM64
  - macOS AMD64 and ARM64
  - Windows AMD64
- Backend plugin support with executable configuration
- Screenshot (GIF) showing plugin in action

### Changed
- Plugin type changed from frontend-only to backend + panel plugin
- Webpack build configuration updated to include backend binaries
- Removed developer-specific badges and jargon from README
- Simplified README by removing development installation section
- Updated API implementation reference to point to bundled test-api

### Fixed
- All Grafana plugin validator warnings resolved
- Binary executable permissions set correctly for all platforms

## [1.0.0] - 2025-10-22

### Added
- **Initial Release**: First stable version of the Hover Tracker Panel
- **Real-time hover tracking**: Captures hover events across Grafana dashboard panels
- **Log correlation**: Automatically correlates metrics with related log entries
- **API integration**: Sends hover data to configurable API endpoints for log analysis
- **Customizable display**: Configurable log limits, truncation, and time windows
- **Modern UI**: Clean, responsive interface with expandable log entries
- **Performance optimization**: Built-in limits for log count and message length
- **Change indicators**: Color-coded log groups based on relative change metrics
- **Zero-configuration viewing**: Just hover - no clicking required
- **Smart time windows**: Queries logs within configurable time windows around hover points
- **API authentication**: Optional Bearer token support for API endpoints
- **Error handling**: Graceful handling of API errors and network issues
- **Responsive design**: Works across different panel sizes and screen resolutions

### Technical Details
- **Grafana compatibility**: Requires Grafana 9.0.0 or higher
- **Node.js requirement**: Built with Node.js 20+ support
- **TypeScript**: Fully typed with comprehensive type definitions
- **React**: Built with React 18.2.0 and modern hooks
- **Grafana UI**: Uses official Grafana UI components and theming
- **Webpack**: Optimized build process with code splitting
- **ESLint**: Code quality enforcement with Grafana's ESLint config
- **Testing**: Jest test framework setup for unit testing

### API Specification
- **Request format**: POST requests with JSON payload containing org, dashboard, panel_title, metric_name, start_time, end_time
- **Response format**: JSON response with log_groups array containing representative_logs and relative_change
- **Authentication**: Optional Bearer token via Authorization header
- **Error handling**: Proper HTTP status codes and error responses
- **CORS support**: Designed to work with cross-origin API requests

### Configuration Options
- **API Endpoint**: Configurable URL for log analysis API
- **API Key**: Optional authentication token
- **Time Window**: Configurable time range for log queries (default: 1 hour)
- **Max Logs**: Maximum number of log entries to display (default: 500)
- **Max Log Length**: Maximum length of individual log entries (default: 10,000)
- **Log Truncate Length**: Character threshold for expandable logs (default: 120)

### Documentation
- **Comprehensive README**: Installation, configuration, and usage instructions
- **API Specification**: Detailed API documentation for server team integration
- **Version Management**: Automated version management scripts and documentation
- **Troubleshooting Guide**: Common issues and solutions
- **Example Implementations**: Sample API implementations in Node.js and Python

### Files Structure
```
src/
├── HoverTrackerPanel.tsx    # Main panel component
├── HoverTrackerEditor.tsx   # Panel configuration editor
├── types.ts                 # TypeScript type definitions
├── plugin.json             # Grafana plugin metadata
└── img/
    └── logo.svg            # Plugin logo

dist/                       # Built plugin files
scripts/
├── version.js             # Version management script
config.json                # Plugin configuration
CHANGELOG.md              # This file
README.md                 # User documentation
API_SPECIFICATION.md      # API documentation
VERSION_MANAGEMENT.md     # Version management guide
```

### Dependencies
- **@grafana/data**: 10.1.5 - Grafana data layer
- **@grafana/runtime**: 10.1.5 - Grafana runtime utilities
- **@grafana/ui**: 10.1.5 - Grafana UI components
- **react**: 18.2.0 - React framework
- **react-dom**: 18.2.0 - React DOM rendering
- **@emotion/css**: ^11.10.6 - CSS-in-JS styling

### Development Tools
- **TypeScript**: 5.0.4 - Type checking and compilation
- **Webpack**: ^5.86.0 - Module bundling
- **ESLint**: Code linting and formatting
- **Jest**: ^29.5.0 - Testing framework
- **Prettier**: ^2.8.7 - Code formatting

### Known Issues
- None in initial release

### Breaking Changes
- None in initial release

### Migration Notes
- This is the initial release, so no migration is required

---

## Release Notes Format

Each release includes:
- **Version number** following semantic versioning (MAJOR.MINOR.PATCH)
- **Release date** in YYYY-MM-DD format
- **Change categories**:
  - **Added**: New features
  - **Changed**: Changes to existing functionality
  - **Deprecated**: Soon-to-be removed features
  - **Removed**: Removed features
  - **Fixed**: Bug fixes
  - **Security**: Security improvements
- **Technical details** for major releases
- **Migration notes** for breaking changes

## Contributing

When contributing to this project, please update this changelog by adding your changes under the appropriate category in the [Unreleased] section. When creating a new release, move the [Unreleased] content to a new version section and update the date.

## Links

- [GitHub Repository](https://github.com/StandardRunbook/grafana-hover-plugin)
- [Plugin Documentation](https://github.com/StandardRunbook/grafana-hover-plugin#readme)
- [API Specification](https://github.com/StandardRunbook/grafana-hover-plugin/blob/main/API_SPECIFICATION.md)
- [Version Management Guide](https://github.com/StandardRunbook/grafana-hover-plugin/blob/main/VERSION_MANAGEMENT.md)
