import { defineNetworkFixture, type NetworkFixture } from "@msw/playwright";
import { expect } from "@playwright/test";
import { HttpNetworkFrame } from "msw/experimental";
import { createBdd, test as baseTest } from "playwright-bdd";
import { defaultHandlers } from "../mocks/handlers/defaultHandlers";

const API_PATH_PREFIX = `${process.env.FRONT_BACKEND_URL}/`;

export const test = baseTest.extend<{ network: NetworkFixture }>({
  network: [
    async ({ context }, use) => {
      const unmockedApiRequests: string[] = [];
      const network = defineNetworkFixture({
        context,
        handlers: defaultHandlers,
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
