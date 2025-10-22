#!/usr/bin/env node

const fs = require('fs');
const path = require('path');

// Get version from command line argument
const newVersion = process.argv[2];

if (!newVersion) {
  console.error('Usage: node scripts/version.js <new-version>');
  console.error('Example: node scripts/version.js 1.1.0');
  process.exit(1);
}

// Validate version format (semantic versioning)
const versionRegex = /^\d+\.\d+\.\d+$/;
if (!versionRegex.test(newVersion)) {
  console.error('Error: Version must be in format X.Y.Z (e.g., 1.1.0)');
  process.exit(1);
}

console.log(`Updating version to ${newVersion}...`);

// Files to update
const filesToUpdate = [
  {
    path: 'package.json',
    update: (content) => {
      const pkg = JSON.parse(content);
      pkg.version = newVersion;
      return JSON.stringify(pkg, null, 2) + '\n';
    }
  },
  {
    path: 'src/plugin.json',
    update: (content) => {
      const plugin = JSON.parse(content);
      plugin.info.version = newVersion;
      plugin.info.updated = new Date().toISOString().split('T')[0];
      return JSON.stringify(plugin, null, 2) + '\n';
    }
  },
  {
    path: 'config.json',
    update: (content) => {
      const config = JSON.parse(content);
      config.plugin.version = newVersion;
      config.build.version = newVersion;
      config.build.buildDate = new Date().toISOString();
      config.files.plugin = `hover-hover-panel-${newVersion}.zip`;
      return JSON.stringify(config, null, 2) + '\n';
    }
  }
];

// Update each file
let updatedFiles = 0;
filesToUpdate.forEach(({ path: filePath, update }) => {
  try {
    const fullPath = path.join(__dirname, '..', filePath);
    
    if (!fs.existsSync(fullPath)) {
      console.warn(`Warning: File ${filePath} not found, skipping...`);
      return;
    }

    const content = fs.readFileSync(fullPath, 'utf8');
    const updatedContent = update(content);
    
    fs.writeFileSync(fullPath, updatedContent);
    console.log(`✓ Updated ${filePath}`);
    updatedFiles++;
  } catch (error) {
    console.error(`Error updating ${filePath}:`, error.message);
  }
});

console.log(`\nVersion update complete! Updated ${updatedFiles} files.`);
console.log(`\nNext steps:`);
console.log(`1. Build the plugin: npm run build`);
console.log(`2. Create zip file: zip -r hover-hover-panel-${newVersion}.zip dist/`);
console.log(`3. Create GitHub release: gh release create v${newVersion} hover-hover-panel-${newVersion}.zip`);
