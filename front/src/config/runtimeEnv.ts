export type AppEnv = {
  FRONT_MOCK_ENABLE: "true" | "false";
  FRONT_BACKEND_URL: string;
  FRONT_DOMAIN: string;
  MODE?: string;
};

const APP_ENV_KEY = "__APP_ENV__" as const;

type RuntimeEnvSource = Partial<Record<string, string | undefined>> & {
  MODE?: string;
};

declare global {
  interface Window {
    [key: string]: unknown;
  }
}

const getSourceEnv = (): RuntimeEnvSource => {
  if (typeof window !== "undefined") {
    const injected = (window as Record<string, unknown>)[APP_ENV_KEY] as
      | RuntimeEnvSource
      | undefined;

    if (injected) {
      return injected;
    }
  }

  if (
    typeof import.meta !== "undefined" &&
    typeof import.meta.env !== "undefined"
  ) {
    return import.meta.env as RuntimeEnvSource;
  }

  if (typeof process !== "undefined" && process.env) {
    return process.env as RuntimeEnvSource;
  }

  return {};
};

export const getRuntimeEnv = (): AppEnv => {
  const sourceEnv = getSourceEnv();

  return {
    FRONT_MOCK_ENABLE: (sourceEnv.FRONT_MOCK_ENABLE ?? "false") as AppEnv["FRONT_MOCK_ENABLE"],
    FRONT_BACKEND_URL: sourceEnv.FRONT_BACKEND_URL ?? "/mocked-api",
    FRONT_DOMAIN: sourceEnv.FRONT_DOMAIN ?? "http://localhost:3000",
    MODE: sourceEnv.MODE ?? "development",
  };
};

export const getRuntimeMode = (): string => {
  const sourceEnv = getSourceEnv();

  return sourceEnv.MODE ?? process.env?.NODE_ENV ?? "development";
};
