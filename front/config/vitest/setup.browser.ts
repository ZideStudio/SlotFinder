// Global setup for Vitest browser tests to load CSS tokens, resets, and global styles

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
