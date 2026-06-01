import { playwright } from "@vitest/browser-playwright";
import { defineConfig } from "vitest/config";
import { getBaseConfig } from "./base";

export default defineConfig(({ mode }) => {
  const base = getBaseConfig(mode);
  return {
    ...base,
    test: {
      ...base.test,
      name: "browser",
      root: new URL("../../", import.meta.url).pathname,
      include: ["src/**/*.browser.test.[jt]sx"],
      setupFiles: ["config/vitest/setup.browser.ts"],
      browser: {
        enabled: true,
        provider: playwright(),
        instances: [
          { browser: "chromium", viewport: { width: 1920, height: 1080 } },
        ],
        api: {
          host: "0.0.0.0",
          allowExec: true,
        },
      },
    },
  };
});
