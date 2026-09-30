import { delay } from "msw";

const isVitestRuntime =
  typeof process !== "undefined" && process.env.VITEST === "true";
const isNonDevelopmentMode = import.meta.env?.MODE !== "development";

export const applyMockDelay = async () => {
  if (isVitestRuntime || isNonDevelopmentMode) {
    return;
  }

  await delay(100);
};
