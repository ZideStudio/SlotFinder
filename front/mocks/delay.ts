import { getRuntimeMode } from "@Front/config/runtimeEnv";
import { delay } from "msw/utils/delay";

const isVitestRuntime =
  typeof process !== "undefined" && process.env.VITEST === "true";
const isNonDevelopmentMode = getRuntimeMode() !== "development";

export const applyMockDelay = async () => {
  if (isVitestRuntime || isNonDevelopmentMode) {
    return;
  }

  await delay(100);
};
