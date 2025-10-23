import type { Configuration } from "webpack";
import path from "path";
import { fileURLToPath } from "url";
import ReplaceInFileWebpackPlugin from "replace-in-file-webpack-plugin";
import CopyWebpackPlugin from "copy-webpack-plugin";
import ForkTsCheckerWebpackPlugin from "fork-ts-checker-webpack-plugin";
import ESLintPlugin from "eslint-webpack-plugin";
import LiveReloadPlugin from "webpack-livereload-plugin";
import type { Compiler } from "webpack";
import fs from "fs";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

/**
 * Webpack plugin to fix source maps by loading source files directly from disk
 * and replacing the sourcesContent with the exact file content
 */
class FixSourceMapPlugin {
  apply(compiler: Compiler) {
    compiler.hooks.emit.tapAsync("FixSourceMapPlugin", (compilation, callback) => {
      Object.keys(compilation.assets).forEach((filename) => {
        if (filename.endsWith(".map")) {
          const asset = compilation.assets[filename];
          const source = asset.source().toString();

          try {
            const sourceMap = JSON.parse(source);

            if (sourceMap.sourcesContent && Array.isArray(sourceMap.sourcesContent)) {
              // Replace sourcesContent with actual file content from disk
              sourceMap.sourcesContent = sourceMap.sourcesContent.map((content: string | null, index: number) => {
                const sourcePath = sourceMap.sources[index];

                // Only process ./src/* files
                if (sourcePath && sourcePath.startsWith('./src/')) {
                  const fullPath = path.resolve(__dirname, "../..", sourcePath);
                  try {
                    // Read the actual file content from disk
                    const fileContent = fs.readFileSync(fullPath, 'utf-8');
                    console.log(`Loaded ${sourcePath} from disk (${fileContent.length} bytes)`);
                    return fileContent;
                  } catch (e) {
                    console.error(`Failed to read ${fullPath}:`, e);
                    return content;
                  }
                }

                // For non-source files (webpack runtime, externals), set to null
                return null;
              });

              const newSource = JSON.stringify(sourceMap);
              compilation.assets[filename] = {
                source: () => newSource,
                size: () => newSource.length,
              } as any;
            }
          } catch (e) {
            console.error(`Failed to process source map ${filename}:`, e);
          }
        }
      });

      callback();
    });
  }
}

const config = (env: any): Configuration => {
  const isProduction = env.production;

  const baseConfig: Configuration = {
    mode: isProduction ? "production" : "development",
    target: "web",
    entry: path.resolve(__dirname, "../../src/module.ts"),
    devtool: isProduction ? "source-map" : "eval-source-map",
    output: {
      path: path.resolve(__dirname, "../../dist"),
      filename: "module.js",
      library: {
        type: "amd",
      },
      clean: true,
      devtoolModuleFilenameTemplate: (info: any) => {
        // Use relative paths starting with ./ for source maps
        const rel = path.relative(
          path.resolve(__dirname, "../.."),
          info.absoluteResourcePath
        );
        return `./${rel.replace(/\\/g, "/")}`;
      },
    },

    externals: [
      "lodash",
      "jquery",
      "moment",
      "slate",
      "emotion",
      "@emotion/react",
      "@emotion/css",
      "prismjs",
      "slate-plain-serializer",
      "@grafana/slate-react",
      "react",
      "react-dom",
      "react-redux",
      "redux",
      "rxjs",
      "react-router",
      "react-router-dom",
      "d3",
      "angular",
      "@grafana/ui",
      "@grafana/runtime",
      "@grafana/data",
      /^@grafana\/.*/,
    ],

    plugins: [
      new CopyWebpackPlugin({
        patterns: [
          { from: "src/plugin.json", to: "." },
          { from: "src/img", to: "./img", noErrorOnMissing: true },
          { from: "README.md", to: "." },
          { from: "CHANGELOG.md", to: ".", noErrorOnMissing: true },
          { from: "LICENSE", to: ".", noErrorOnMissing: true },
          { from: "bin/grafana-plugin-api-darwin-amd64", to: "gpx_hover-hover-panel_darwin_amd64", toType: "file", noErrorOnMissing: true },
          { from: "bin/grafana-plugin-api-darwin-arm64", to: "gpx_hover-hover-panel_darwin_arm64", toType: "file", noErrorOnMissing: true },
          { from: "bin/grafana-plugin-api-linux-amd64", to: "gpx_hover-hover-panel_linux_amd64", toType: "file", noErrorOnMissing: true },
          { from: "bin/grafana-plugin-api-linux-arm64", to: "gpx_hover-hover-panel_linux_arm64", toType: "file", noErrorOnMissing: true },
          { from: "bin/grafana-plugin-api-windows-amd64.exe", to: "gpx_hover-hover-panel_windows_amd64.exe", toType: "file", noErrorOnMissing: true },
        ],
      }),

      new ReplaceInFileWebpackPlugin([
        {
          dir: "dist",
          files: ["plugin.json"],
          rules: [
            {
              search: "%VERSION%",
              replace: "1.0.0",
            },
            {
              search: "%TODAY%",
              replace: new Date().toISOString().substring(0, 10),
            },
          ],
        },
      ]) as any,

      new ForkTsCheckerWebpackPlugin({
        async: Boolean(isProduction),
        typescript: {
          configFile: path.resolve(__dirname, "../../tsconfig.json"),
        },
      }),

      new ESLintPlugin({
        extensions: [".ts", ".tsx"],
        lintDirtyModulesOnly: !isProduction,
      }),

      new FixSourceMapPlugin(),
    ],

    resolve: {
      extensions: [".js", ".jsx", ".ts", ".tsx"],
      modules: [path.resolve(__dirname, "../../src"), "node_modules"],
    },

    module: {
      rules: [
        {
          test: /\.(ts|tsx)$/,
          exclude: /node_modules/,
          use: {
            loader: "ts-loader",
            options: {
              transpileOnly: false,
              configFile: path.resolve(__dirname, "../../tsconfig.json"),
            },
          },
        },
        {
          test: /\.css$/,
          use: ["style-loader", "css-loader"],
        },
        {
          test: /\.(png|jpe?g|gif|svg)$/,
          type: "asset/resource",
        },
      ],
    },
  };

  if (!isProduction) {
    baseConfig.plugins!.push(new LiveReloadPlugin() as any);
  }

  return baseConfig;
};

export default config;
