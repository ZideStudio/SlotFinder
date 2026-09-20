import { playwright } from "@vitest/browser-playwright";
import { defineConfig } from "vitest/config";
import { getBaseConfig } from "./base";

export default defineConfig(({ mode }) => {
  const base = getBaseConfig(mode);
  return {
    ...base,
    publicDir: new URL("../../public", import.meta.url).pathname,
    test: {
      ...base.test,
      name: "browser",
      root: new URL("../../", import.meta.url).pathname,
      include: ["src/**/*.browser.test.[jt]sx"],
      setupFiles: ["config/vitest/setup.browser.ts"],
      api: {
        host: "0.0.0.0",
        allowExec: true,
      },
      browser: {
        enabled: true,
        headless: true,
        provider: playwright(),
        instances: [
          {
            browser: "chromium",
            viewport: { width: 1920, height: 1080 },
          },
        ],
      },
    },
  };
});
