#!/bin/bash

# Grafana Plugin Release Script
# Usage: ./scripts/release.sh <version>
# Example: ./scripts/release.sh 1.0.8

set -e

VERSION=$1
PLUGIN_ID="hover-hover-panel"
PLUGIN_NAME="hover-hover-panel"
REPO="StandardRunbook/grafana-hover-plugin"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

if [ -z "$VERSION" ]; then
  echo -e "${RED}Error: Version number required${NC}"
  echo "Usage: $0 <version>"
  echo "Example: $0 1.0.8"
  exit 1
fi

if [ -z "$GRAFANA_ACCESS_POLICY_TOKEN" ]; then
  echo -e "${RED}Error: GRAFANA_ACCESS_POLICY_TOKEN environment variable not set${NC}"
  echo "Export your Grafana signing token first:"
  echo "export GRAFANA_ACCESS_POLICY_TOKEN=your_token_here"
  exit 1
fi

echo -e "${GREEN}🚀 Starting release process for v${VERSION}${NC}"
echo

# Step 1: Update version in package.json and plugin.json
echo -e "${YELLOW}📝 Updating version to ${VERSION}...${NC}"
sed -i.bak "s/\"version\": \".*\"/\"version\": \"${VERSION}\"/" package.json
sed -i.bak "s/\"version\": \".*\"/\"version\": \"${VERSION}\"/" src/plugin.json
rm -f package.json.bak src/plugin.json.bak

# Step 2: Build the plugin
echo -e "${YELLOW}🔨 Building plugin...${NC}"
pnpm run build

# Step 3: Set executable permissions on binaries
echo -e "${YELLOW}🔐 Setting executable permissions...${NC}"
chmod +x dist/gpx_${PLUGIN_NAME}_*

# Step 4: Copy Go manifest before signing
echo -e "${YELLOW}📋 Copying Go manifest...${NC}"
cp pkg/go_plugin_build_manifest dist/

# Step 5: Sign the plugin
echo -e "${YELLOW}✍️  Signing plugin...${NC}"
npx --yes @grafana/sign-plugin@latest --rootUrls http://localhost:3000

# Step 6: Create plugin package
echo -e "${YELLOW}📦 Creating plugin package...${NC}"
rm -rf tmp-package
mkdir -p tmp-package/${PLUGIN_NAME}
cp -r dist/* tmp-package/${PLUGIN_NAME}/
cd tmp-package
zip -r ../${PLUGIN_NAME}-${VERSION}.zip ${PLUGIN_NAME}
cd ..
rm -rf tmp-package

# Step 7: Generate checksums
echo -e "${YELLOW}🔒 Generating checksums...${NC}"
MD5=$(md5sum ${PLUGIN_NAME}-${VERSION}.zip | awk '{print $1}')
SHA256=$(shasum -a 256 ${PLUGIN_NAME}-${VERSION}.zip | awk '{print $1}')
echo "$MD5" > ${PLUGIN_NAME}-${VERSION}.zip.md5
echo "$SHA256" > ${PLUGIN_NAME}-${VERSION}.zip.sha256

# Step 8: Commit changes
echo -e "${YELLOW}💾 Committing changes...${NC}"
git add -A
git commit -m "Release v${VERSION}

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>" || echo "No changes to commit"

# Step 9: Push to GitHub
echo -e "${YELLOW}⬆️  Pushing to GitHub...${NC}"
git push

# Step 10: Create GitHub release
echo -e "${YELLOW}🎉 Creating GitHub release...${NC}"
gh release create v${VERSION} \
  --target main \
  ${PLUGIN_NAME}-${VERSION}.zip \
  ${PLUGIN_NAME}-${VERSION}.zip.md5 \
  ${PLUGIN_NAME}-${VERSION}.zip.sha256 \
  --title "Hover Panel v${VERSION}" \
  --notes "## Hover Panel v${VERSION}

Grafana plugin for automatic log correlation on hover.

### 📦 Installation

\`\`\`bash
unzip ${PLUGIN_NAME}-${VERSION}.zip -d /var/lib/grafana/plugins/
\`\`\`

### 🔐 Checksums

- **MD5**: \`${MD5}\`
- **SHA256**: \`${SHA256}\`

### 📚 Documentation

- [README](https://github.com/${REPO}#readme)
- [Testing Guide](https://github.com/${REPO}/blob/main/TESTING_GUIDE.md)

---

🤖 Generated with [Claude Code](https://claude.com/claude-code)"

echo
echo -e "${GREEN}✅ Release v${VERSION} completed successfully!${NC}"
echo
echo -e "${GREEN}📋 Submission Information:${NC}"
echo
echo "Plugin ID: ${PLUGIN_ID}"
echo
echo "Download URL:"
echo "https://github.com/${REPO}/releases/download/v${VERSION}/${PLUGIN_NAME}-${VERSION}.zip"
echo
echo "MD5 Checksum:"
echo "${MD5}"
echo
echo "SHA256 Checksum:"
echo "${SHA256}"
echo
echo "Source Repository:"
echo "https://github.com/${REPO}"
echo
