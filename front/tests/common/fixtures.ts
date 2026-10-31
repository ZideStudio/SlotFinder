import { defineNetworkFixture, type NetworkFixture } from "@msw/playwright";
import { expect } from "@playwright/test";
import { HttpNetworkFrame } from "msw/experimental";
import { test as baseTest, createBdd } from "playwright-bdd";
import { getAuthStatus401 } from "../../mocks/handlers/authStatusHandlers";
import { defaultHandlers } from "../../mocks/handlers/defaultHandlers";

const API_PATH_PREFIX = `${process.env.FRONT_BACKEND_URL}/`;

export const test = baseTest.extend<{ network: NetworkFixture }>({
  network: [
    async ({ context, $tags }, use) => {
      const unmockedApiRequests: string[] = [];
      const network = defineNetworkFixture({
        context,
        // First match wins, so the 401 overrides the signed-in baseline.
        handlers: $tags.includes("@unauthenticated")
          ? [getAuthStatus401, ...defaultHandlers]
          : defaultHandlers,
        // `defaults.error()` only logs, so the frame is rejected explicitly.
        onUnhandledFrame: ({ frame, defaults }) => {
          if (!(frame instanceof HttpNetworkFrame)) {
            return;
          }

          const { request } = frame.data;

          if (!new URL(request.url).pathname.startsWith(API_PATH_PREFIX)) {
            return;
          }

          unmockedApiRequests.push(`${request.method} ${request.url}`);
          defaults.error();
          throw new Error(`Unmocked API request: ${request.url}`);
        },
      });

      await network.enable();
      await use(network);
      await network.disable();

      // The app may swallow the aborted request, so fail the test here too.
      expect(unmockedApiRequests, "Unmocked API requests").toEqual([]);
    },
    { auto: true },
  ],
});

export const { Given, When, Then } = createBdd(test);
