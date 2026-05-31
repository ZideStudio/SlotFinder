// Global setup for Vitest browser tests: loads global styles and i18n, and manages the MSW service worker lifecycle.

// oxlint-disable vitest/require-top-level-describe
import "@Front/assets/css/index";
import "@Front/i18n/index";
import { worker } from "@Mocks/browser";

beforeAll(async () => {
  await worker.start({ onUnhandledRequest: "bypass" });
});

afterEach(() => {
  worker.resetHandlers();
});

afterAll(() => {
  worker.stop();
});
