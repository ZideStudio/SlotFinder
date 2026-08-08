import { defineConfig, loadEnv } from "@rsbuild/core";
import { pluginReact } from "@rsbuild/plugin-react";
import { pluginSass } from "@rsbuild/plugin-sass";
import { pluginSvgr } from "@rsbuild/plugin-svgr";
import { resolve } from "node:path";
import packageJson from "./package.json";

const { publicVars } = loadEnv({
  prefixes: [process.env.ENV_PREFIX ?? "FRONT_"],
});

const publicUrl =
  process.env.PUBLIC_URL ??
  (packageJson as { homepage?: string }).homepage ??
  "/";
const publicPath = new URL(
  publicUrl.endsWith("/") ? publicUrl : `${publicUrl}/`,
  "http://localhost",
).pathname;
const disableSourceMap =
  (process.env.DISABLE_SOURCE_MAP ?? "false") === "true" ? false : "source-map";
export default defineConfig(({ env }) => {
  const isProduction = env === "production";

  return {
    plugins: [
      pluginReact(),
      pluginSvgr({
        query: /.*/,
        svgrOptions: {
          exportType: "default",
        },
      }),
      pluginSass(),
    ],
    source: {
      entry: {
        index: "./src/main.ts",
      },
      define: publicVars,
    },
    resolve: {
      alias: {
        "@Front": resolve(__dirname, "src"),
        "@Mocks": resolve(__dirname, "mocks"),
      },
    },
    output: {
      assetPrefix: publicPath,
      sourceMap: {
        js: isProduction ? disableSourceMap : "cheap-module-source-map",
        css: (process.env.DISABLE_SOURCE_MAP ?? "false") !== "true",
      },
      distPath: {
        root: "build",
      },
      copy: [
        {
          from: "public",
          to: ".",
        },
      ],
    },
    server: {
      host: process.env.HOST ?? "0.0.0.0",
      port: Number.parseInt(process.env.PORT ?? "3000", 10),
      open: (process.env.BROWSER ?? "false") === "true",
    },
    html: {
      template: "./index.html",
    },
  };
});
