import { createBdd, test as baseTest } from "playwright-bdd";
import { defineNetworkFixture, type NetworkFixture } from "@msw/playwright";
import {
  getAccountAvatar200,
  getAccountMe200,
} from "../mocks/handlers/accountHandlers";
import { getAuthStatus401 } from "../mocks/handlers/authStatusHandlers";

const APP_ENV = {
  FRONT_MOCK_ENABLE: "false",
  FRONT_BACKEND_URL: "/mocked-api",
  FRONT_DOMAIN: "http://127.0.0.1:3004",
  MODE: "development",
} as const;

export const test = baseTest.extend<{ network: NetworkFixture }>({
  network: [
    async ({ context }, use) => {
      const network = defineNetworkFixture({
        context,
        handlers: [getAuthStatus401, getAccountMe200, getAccountAvatar200],
      });

      await network.enable();
      await use(network);
      await network.disable();
    },
    { auto: true },
  ],
});

const { Before, Given, When, Then } = createBdd(test);

Before(async ({ page }) => {
  await page.addInitScript((env) => {
    const runtimeKey = "__APP_ENV__" as const;
    (globalThis as Record<string, unknown>)[runtimeKey] = env;
  }, APP_ENV);
});

export { Before, Given, When, Then };
