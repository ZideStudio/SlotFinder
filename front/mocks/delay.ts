import { delay } from "msw";

const isVitestRuntime =
  typeof process !== "undefined" && process.env.VITEST === "true";
const isViteTestMode = import.meta.env?.MODE !== "development";

export const applyMockDelay = async () => {
  if (isVitestRuntime || isViteTestMode) {
    return;
  }

  await delay(100);
};
