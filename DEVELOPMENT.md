# Development Guide

This document contains developer and contributor information for the Grafana Hover Tracker Panel plugin.

## Prerequisites

- **Node.js**: 20 or higher
- **pnpm**: 8 or higher (recommended package manager)
- **Grafana**: 9.0.0 or higher

## Development Setup

```bash
# Clone the repository
git clone https://github.com/StandardRunbook/grafana-hover-plugin.git
cd grafana-hover-plugin

# Install dependencies
pnpm install

# Start development server
pnpm run dev
```

## Available Scripts

### Development
```bash
pnpm run dev          # Start development mode with watch
pnpm run build        # Build the plugin for production
pnpm run typecheck    # Run TypeScript type checking
```

### Testing
```bash
pnpm run test         # Run tests in watch mode
pnpm run test:ci      # Run tests for CI
pnpm run e2e          # Run end-to-end tests
```

### Code Quality
```bash
pnpm run lint         # Run ESLint
pnpm run lint:fix     # Fix ESLint issues automatically
```

### Version Management
```bash
pnpm run version 1.1.0    # Update version across all files
pnpm run package          # Build and create zip file
pnpm run release 1.1.0   # Update version and create package
```

### Plugin Management
```bash
pnpm run sign            # Sign the plugin for distribution
```

## Project Structure

```
src/
├── HoverTrackerPanel.tsx    # Main panel component
├── HoverTrackerEditor.tsx   # Panel configuration editor
├── types.ts                 # TypeScript type definitions
├── plugin.json             # Grafana plugin metadata
└── img/
    └── logo.svg            # Plugin logo

dist/                       # Built plugin files (generated)
scripts/
├── version.js             # Version management script
config.json                # Plugin configuration
CHANGELOG.md              # Release changelog
README.md                 # User documentation
API_SPECIFICATION.md      # API documentation
DEVELOPMENT.md            # This file
VERSION_MANAGEMENT.md     # Version management guide
```

## Architecture

### Core Components

1. **HoverTrackerPanel**: Main React component that renders the panel UI
2. **HoverTrackerEditor**: Configuration editor for panel settings
3. **Event Handling**: Subscribes to Grafana hover events and processes them
4. **API Integration**: Sends hover data to configured API endpoints

### Key Features

- **Real-time hover tracking**: Uses Grafana's event bus to capture hover events
- **API communication**: Sends structured data to external log analysis APIs
- **Responsive UI**: Adapts to different panel sizes and screen resolutions
- **Error handling**: Graceful handling of API errors and network issues

## API Integration

The plugin sends POST requests to configured API endpoints with the following structure:

```json
{
  "org": "1",
  "dashboard": "Dashboard Name",
  "panel_title": "Panel Title",
  "metric_name": "Series/Metric Name",
  "start_time": "2024-01-15T10:30:00.000Z",
  "end_time": "2024-01-15T11:30:00.000Z"
}
```

Expected response format:

```json
{
  "log_groups": [
    {
      "representative_logs": ["Log entry 1", "Log entry 2"],
      "relative_change": 15.5
    }
  ]
}
```

## Testing

### Unit Tests
- Located in `__tests__/` directories
- Use Jest and React Testing Library
- Run with `pnpm run test`

### End-to-End Tests
- Use Grafana's E2E testing framework
- Test plugin functionality in real Grafana environment
- Run with `pnpm run e2e`

### Manual Testing
- Use the provided Docker Compose setup
- Includes mock API server and sample data
- Run with `docker-compose -f docker-compose.test.yml up`

## Building and Packaging

### Development Build
```bash
pnpm run build
```

### Production Package
```bash
pnpm run package
```

This creates a zip file ready for distribution.

### Plugin Signing
```bash
export GRAFANA_ACCESS_POLICY_TOKEN=<your-token>
pnpm run sign
```

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Make your changes
4. Add tests if applicable
5. Run the test suite: `pnpm run test:ci`
6. Run linting: `pnpm run lint`
7. Build the plugin: `pnpm run build`
8. Commit your changes: `git commit -m 'Add amazing feature'`
9. Push to your branch: `git push origin feature/amazing-feature`
10. Open a Pull Request

## Code Style

- Follow the existing code style
- Use TypeScript for all new code
- Add proper type definitions
- Include JSDoc comments for public APIs
- Use meaningful variable and function names

## Debugging

### Browser DevTools
- Use React DevTools for component debugging
- Check Network tab for API requests
- Monitor Console for error messages

### Grafana Logs
- Check Grafana server logs for plugin errors
- Enable debug logging in Grafana configuration

### Plugin Validator
```bash
pnpm exec plugin-validator dist/
```

## Release Process

1. Update version: `pnpm run version 1.1.0`
2. Update CHANGELOG.md with new features/fixes
3. Build and test: `pnpm run build && pnpm run test:ci`
4. Create release: `gh release create v1.1.0`
5. GitHub Actions will automatically build and package the plugin

## Troubleshooting

### Common Issues

**Plugin not loading:**
- Check Grafana logs for errors
- Verify plugin files are in correct directory
- Ensure plugin is signed or allow unsigned plugins

**API requests failing:**
- Check CORS settings on your API
- Verify API endpoint URL is correct
- Check browser network tab for error details

**Build errors:**
- Ensure all dependencies are installed: `pnpm install`
- Check TypeScript errors: `pnpm run typecheck`
- Verify Node.js version compatibility

## Resources

- [Grafana Plugin Development](https://grafana.com/developers/plugin-tools/)
- [React Documentation](https://reactjs.org/docs/)
- [TypeScript Handbook](https://www.typescriptlang.org/docs/)
- [Grafana UI Components](https://developers.grafana.com/ui/latest/)
