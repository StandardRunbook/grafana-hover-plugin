# Version Management

This document describes how to manage versions for the Grafana Hover Tracker Panel plugin.

## Files That Track Version

The following files contain version information that need to be kept in sync:

1. **`package.json`** - Main version for npm/Node.js
2. **`src/plugin.json`** - Grafana plugin metadata
3. **`config.json`** - Plugin configuration and metadata

## Version Management Script

A Node.js script is provided to automatically update versions across all files:

### Usage

```bash
# Update to a new version
npm run version 1.1.0

# Or run the script directly
node scripts/version.js 1.1.0
```

### What the Script Does

1. Validates the version format (must be X.Y.Z)
2. Updates version in `package.json`
3. Updates version and build date in `src/plugin.json`
4. Updates version, build date, and zip filename in `config.json`
5. Reports which files were updated

## Available Scripts

### Version Management
```bash
npm run version <new-version>    # Update version across all files
npm run package                  # Build and create zip file
npm run release <new-version>    # Update version and create package
```

### Development
```bash
npm run build                    # Build the plugin
npm run dev                      # Development mode with watch
npm run test                     # Run tests
npm run lint                     # Lint code
```

## Release Process

1. **Update Version**:
   ```bash
   npm run version 1.1.0
   ```

2. **Build and Package**:
   ```bash
   npm run package
   ```

3. **Create GitHub Release**:
   ```bash
   gh release create v1.1.0 hover-hover-panel-1.1.0.zip --title "Release v1.1.0" --notes "Release notes here"
   ```

## Configuration File

The `config.json` file contains:

- **Plugin metadata**: ID, name, version, author info
- **Build information**: Version, build date, requirements
- **API configuration**: Endpoint, method, authentication
- **Default settings**: Time windows, log limits
- **Repository info**: Git URL, releases URL
- **File references**: Plugin zip, API spec, documentation

## Version Format

Versions follow [Semantic Versioning](https://semver.org/):
- **MAJOR**: Incompatible API changes
- **MINOR**: New functionality in a backwards compatible manner
- **PATCH**: Backwards compatible bug fixes

Examples: `1.0.0`, `1.1.0`, `1.1.1`, `2.0.0`

## Troubleshooting

### Version Script Errors

- **"Version must be in format X.Y.Z"**: Use semantic versioning format
- **"File not found"**: Ensure you're running from the project root
- **"Error updating file"**: Check file permissions and JSON validity

### Manual Version Updates

If the script fails, you can manually update versions in:
1. `package.json` - `version` field
2. `src/plugin.json` - `info.version` and `info.updated` fields
3. `config.json` - `plugin.version`, `build.version`, `build.buildDate`, and `files.plugin` fields

## Best Practices

1. **Always use the version script** to ensure consistency
2. **Test after version updates** to ensure everything works
3. **Commit version changes** before creating releases
4. **Update changelog** when releasing new versions
5. **Tag releases** in git for better tracking
