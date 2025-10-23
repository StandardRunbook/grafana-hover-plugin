declare module 'replace-in-file-webpack-plugin' {
  import { Compiler } from 'webpack';

  interface ReplaceRule {
    search: string | RegExp;
    replace: string;
  }

  interface ReplaceInFileOptions {
    dir: string;
    files: string[];
    rules: ReplaceRule[];
  }

  class ReplaceInFileWebpackPlugin {
    constructor(options: ReplaceInFileOptions[]);
    apply(compiler: Compiler): void;
  }

  export = ReplaceInFileWebpackPlugin;
}
