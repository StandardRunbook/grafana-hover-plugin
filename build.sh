#!/bin/bash
set -e

echo "Building with Grafana Plugin SDK..."

# Step 1: Clean
rm -rf dist bin

# Step 2: Build Go binaries with mage
mage buildAll

# Step 3: Copy binaries to bin/ for webpack
mkdir -p bin
for src in dist/gpx_*; do
  dst="bin/$(basename "$src" | sed 's/gpx_hover-hover-panel_/grafana-plugin-api-/')"
  cp "$src" "$dst"
done

# Step 4: Build frontend with webpack
pnpm run build

# Step 5: Copy manifest back to dist after webpack
cp go_plugin_build_manifest dist/

echo "✓ Build complete!"
