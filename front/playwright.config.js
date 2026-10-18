import { defineConfig } from "@playwright/test";
import dotenv from "dotenv";
import { defineBddConfig } from "playwright-bdd";

dotenv.config({ path: ".env.mock" });

const testDir = defineBddConfig({
  features: "tests/**/*.feature",
  steps: "tests/**/*.steps.ts",
  outputDir: ".features-gen",
});

export default defineConfig({
  testDir,
  reporter: [["html", { outputFolder: "playwright-report" }]],
  fullyParallel: false,
  use: {
    baseURL: "http://127.0.0.1:3004",
    trace: "on-first-retry",
    screenshot: "only-on-failure",
    video: "retain-on-failure",
    headless: true,
  },
  webServer: {
    command: "npm run start -- --host 0.0.0.0 --port 3004",
    env: {
      FRONT_MOCK_ENABLE: "false",
      FRONT_BACKEND_URL: "/mocked-api",
      FRONT_DOMAIN: "https://localhost:3004",
    },
    url: "http://127.0.0.1:3004",
    reuseExistingServer: !process.env.CI,
    timeout: 180_000,
  },
  ui: true,
  uiPort: 9323,
  uiHost: "0.0.0.0",
});
