import { delay } from "msw/utils/delay";

const isVitestRuntime =
  typeof process !== "undefined" && process.env.VITEST === "true";
const isNonDevelopmentMode = import.meta.env?.MODE !== "development";

export const applyMockDelay = async () => {
  if (isVitestRuntime || isNonDevelopmentMode) {
    return;
  }

  await delay(100);
};
