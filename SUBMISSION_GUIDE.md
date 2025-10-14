# Grafana Plugin Submission Guide

This guide walks you through submitting the Hover panel plugin to the official Grafana plugin catalog.

## Pre-Submission Checklist

### ✅ Completed
- [x] Plugin ID updated to `hover-hover-panel`
- [x] Author information added (Ankil Patel)
- [x] LICENSE file added (Apache 2.0)
- [x] Logo created (src/img/logo.svg)
- [x] Comprehensive README.md created
- [x] Panel options configurable via UI
- [x] Links added to plugin.json (GitHub, docs, website)

### ⚠️ Still Needed
- [ ] **Create screenshots** (see SCREENSHOTS.md)
  - screenshot-main.png - showing logs on hover
  - screenshot-config.png - showing panel configuration
- [ ] **Push code to public GitHub repository**
  - Repository: https://github.com/ankilp/grafana-hover-tracker-panel
- [ ] **Create GitHub release** with signed plugin ZIP
- [ ] **Prepare API backend documentation** (users need to implement their own)

## Step-by-Step Submission Process

### Step 1: Create Screenshots

Before submitting, you need real screenshots. See [SCREENSHOTS.md](SCREENSHOTS.md) for details.

1. Install the plugin in a Grafana instance
2. Create a demo dashboard with time-series panels
3. Take screenshots showing:
   - The plugin displaying logs when hovering over metrics
   - The configuration panel with all options visible
4. Save as PNG in `src/img/`:
   - `screenshot-main.png`
   - `screenshot-config.png`
5. Rebuild: `pnpm run build`

### Step 2: Push to GitHub

```bash
# Initialize git if not already done
git init
git add .
git commit -m "Initial commit - Hover panel plugin"

# Add remote (update with your actual repo URL)
git remote add origin https://github.com/ankilp/grafana-hover-tracker-panel.git
git branch -M main
git push -u origin main
```

### Step 3: Sign the Plugin

You need a Grafana Cloud account to sign plugins.

#### Create Access Policy Token

1. Sign in to [Grafana Cloud](https://grafana.com)
2. Go to **Administration** → **Cloud Access Policies**
3. Click **Create Access Policy**
4. Name: "Plugin Signing - Hover Panel"
5. Scopes: Select **plugins:write**
6. Click **Add token**
7. Display name: "Hover Panel Signing Token"
8. Click **Create** and copy the token

#### Sign Your Plugin

```bash
# Export the token
export GRAFANA_ACCESS_POLICY_TOKEN=<your-token-here>

# Build the plugin first
pnpm run build

# Sign the plugin
npx @grafana/sign-plugin@latest

# This creates a MANIFEST.txt file in dist/
```

### Step 4: Create GitHub Release

```bash
# Tag the version
git tag -a v1.0.0 -m "Initial release"
git push origin v1.0.0

# Package the signed plugin
cd dist
zip -r ../hover-hover-panel-1.0.0.zip .
cd ..
```

#### Create Release on GitHub

1. Go to your repository on GitHub
2. Click **Releases** → **Create a new release**
3. Tag version: `v1.0.0`
4. Release title: `Hover Panel v1.0.0`
5. Description:
   ```markdown
   ## First Release
   
   Hover panel plugin that displays related logs when hovering over metrics.
   
   ### Features
   - Automatic log correlation on hover
   - Configurable API endpoint
   - Expandable log entries
   - Time window configuration
   - Performance limits
   
   ### Installation
   Download the ZIP file and extract to your Grafana plugins directory.
   ```
6. Upload `hover-hover-panel-1.0.0.zip`
7. Click **Publish release**

### Step 5: Submit to Grafana

1. Go to [Grafana Cloud](https://grafana.com/auth/sign-in)
2. Navigate to **Org Settings** → **My Plugins**
3. Click **Submit New Plugin**
4. Fill out the form:

   **Plugin Archive:**
   - Type: **Single** (one ZIP for all architectures)
   - URL: `https://github.com/ankilp/grafana-hover-tracker-panel/releases/download/v1.0.0/hover-hover-panel-1.0.0.zip`

   **Source Code:**
   - URL: `https://github.com/ankilp/grafana-hover-tracker-panel`

   **Plugin Information:**
   - Plugin ID: `hover-hover-panel`
   - Version: `1.0.0`

5. Click **Submit for Review**

### Step 6: Wait for Review

- Grafana team reviews all submissions
- Review typically takes **1-2 weeks**
- They check:
  - Code quality and security
  - No data collection/privacy violations
  - Proper licensing
  - Working functionality
  - Documentation quality

- You'll receive email notifications about review status
- May receive feedback requesting changes

### Step 7: Address Feedback (if any)

If Grafana requests changes:

1. Make the required changes
2. Commit and push to GitHub
3. Create a new release (e.g., v1.0.1)
4. Sign the new version
5. In Grafana Cloud → My Plugins → Click **Submit Update**
6. Provide new ZIP URL

## Automation with GitHub Actions

To automate signing and releases, create `.github/workflows/release.yml`:

```yaml
name: Release Plugin

on:
  push:
    tags:
      - 'v*'

jobs:
  release:
    runs-on: ubuntu-latest
    
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '18'
      
      - name: Setup pnpm
        uses: pnpm/action-setup@v2
        with:
          version: 8
      
      - name: Install dependencies
        run: pnpm install
      
      - name: Build plugin
        run: pnpm run build
      
      - name: Sign plugin
        run: npx @grafana/sign-plugin@latest
        env:
          GRAFANA_ACCESS_POLICY_TOKEN: ${{ secrets.GRAFANA_ACCESS_POLICY_TOKEN }}
      
      - name: Package plugin
        run: |
          cd dist
          zip -r ../hover-hover-panel-${{ github.ref_name }}.zip .
      
      - name: Create GitHub Release
        uses: softprops/action-gh-release@v1
        with:
          files: hover-hover-panel-${{ github.ref_name }}.zip
          body: |
            ## Hover Panel ${{ github.ref_name }}
            
            Download the ZIP and extract to your Grafana plugins directory.
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

**Setup:**
1. Go to your GitHub repository → Settings → Secrets and variables → Actions
2. Add new repository secret:
   - Name: `GRAFANA_ACCESS_POLICY_TOKEN`
   - Value: Your Grafana access policy token

**Usage:**
```bash
git tag v1.0.0
git push origin v1.0.0
# GitHub Actions automatically builds, signs, and creates release
```

## Important Notes

### API Backend Requirement

This plugin requires users to implement their own log analysis API. Make sure to clearly document:
- Expected request format
- Expected response format
- Example implementations
- CORS requirements

Users cannot use the plugin without this backend component.

### Plugin Naming

The plugin ID `hover-hover-panel` follows Grafana's convention:
- `hover` - organization/company name
- `hover` - plugin name
- `panel` - plugin type

This might look redundant, but it follows the pattern: `<org>-<name>-<type>`

### Signature Levels

Plugins can have different signature levels:
- **Community**: Free, for open source plugins
- **Commercial**: Requires paid subscription
- **Private**: For internal use only

Your plugin will be **Community** signed (free and public).

### Updates

For future updates:
1. Increment version in package.json
2. Tag and push: `git tag v1.0.1 && git push origin v1.0.1`
3. Build, sign, and package
4. In Grafana Cloud → My Plugins → Click **Submit Update**
5. Provide new version ZIP URL

## Support After Publishing

Once published, users may:
- Report issues on GitHub
- Request features
- Submit pull requests
- Ask questions in discussions

Be prepared to:
- Respond to issues
- Maintain documentation
- Release updates for Grafana version compatibility
- Provide example API implementations

## Resources

- [Grafana Plugin Tools](https://grafana.com/developers/plugin-tools/)
- [Plugin Publishing Guide](https://grafana.com/developers/plugin-tools/publish-a-plugin/publish-a-plugin)
- [Plugin Signing Guide](https://grafana.com/developers/plugin-tools/publish-a-plugin/sign-a-plugin)
- [Grafana Plugin Catalog](https://grafana.com/grafana/plugins/)

## Questions?

If you have questions during the submission process:
- Check the [Grafana Community Forums](https://community.grafana.com/)
- Review [Grafana Plugin Documentation](https://grafana.com/docs/grafana/latest/developers/plugins/)
- Email Grafana plugin team (provided after submission)
