declare module 'webpack-livereload-plugin' {
  import { Compiler } from 'webpack';

  interface LiveReloadPluginOptions {
    protocol?: 'http' | 'https';
    port?: number;
    hostname?: string;
    appendScriptTag?: boolean;
    ignore?: RegExp | null;
    delay?: number;
  }

  class LiveReloadPlugin {
    constructor(options?: LiveReloadPluginOptions);
    apply(compiler: Compiler): void;
  }

  export = LiveReloadPlugin;
}
