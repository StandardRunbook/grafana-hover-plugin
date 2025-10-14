# Quick Submission Checklist

Use this as a quick reference when you're ready to submit.

## Before You Start
- [ ] Screenshots created (see SCREENSHOTS.md)
- [ ] Code pushed to GitHub: https://github.com/ankilp/grafana-hover-tracker-panel
- [ ] Grafana Cloud account created

## Signing & Release

```bash
# 1. Get Grafana Cloud Access Token
# Go to: https://grafana.com → Administration → Cloud Access Policies
# Create policy with plugins:write scope

# 2. Export token
export GRAFANA_ACCESS_POLICY_TOKEN=<your-token>

# 3. Build and sign
pnpm run build
npx @grafana/sign-plugin@latest

# 4. Package
cd dist
zip -r ../hover-hover-panel-1.0.0.zip .
cd ..

# 5. Tag and push
git tag v1.0.0
git push origin v1.0.0

# 6. Create GitHub release
# Upload hover-hover-panel-1.0.0.zip
# Get download URL from release page
```

## Submit to Grafana

1. Go to: https://grafana.com/auth/sign-in
2. Navigate to: **Org Settings** → **My Plugins**
3. Click: **Submit New Plugin**
4. Fill form:
   - Architecture: **Single**
   - ZIP URL: `https://github.com/ankilp/grafana-hover-tracker-panel/releases/download/v1.0.0/hover-hover-panel-1.0.0.zip`
   - Source: `https://github.com/ankilp/grafana-hover-tracker-panel`
   - Plugin ID: `hover-hover-panel`
   - Version: `1.0.0`
5. Submit and wait for review (1-2 weeks)

## Files Created

✅ Updated:
- `src/plugin.json` - Plugin metadata with your info
- `README.md` - Comprehensive documentation

✅ Created:
- `LICENSE` - Apache 2.0 license
- `src/img/logo.svg` - Plugin logo
- `SCREENSHOTS.md` - Screenshot guide
- `SUBMISSION_GUIDE.md` - Detailed submission process
- `QUICK_SUBMIT.md` - This checklist

⚠️ Still needed:
- `src/img/screenshot-main.png`
- `src/img/screenshot-config.png`
- GitHub repository with code
- GitHub release with signed ZIP

## Plugin Info

- **ID**: `hover-hover-panel`
- **Name**: Hover
- **Author**: Ankil Patel
- **Company**: Hover
- **License**: Apache 2.0
- **Grafana**: >= 9.0.0
- **Type**: Panel Plugin

## Next Steps After Approval

Once approved:
1. Plugin appears in Grafana catalog
2. Users can install with: `grafana-cli plugins install hover-hover-panel`
3. Monitor GitHub issues for bug reports
4. Release updates following same process
5. Update version in package.json for each release
