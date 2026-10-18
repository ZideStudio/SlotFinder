import { defineConfig } from "orval";

export default defineConfig({
  slotfinder: {
    input: {
      target: "./src/api/orval/swagger.json",
    },
    output: {
      mode: "tags-split",
      target: "./src/api/generated",
      client: "fetch",
      baseUrl: {
        runtime: "getRuntimeEnv().FRONT_BACKEND_URL",
        imports: [{
          name: "getRuntimeEnv",
          importPath: "@Front/config/runtimeEnv",
        }],
      },
      override: {
        fetch: {
          includeHttpResponseReturnType: false,
        },
        mutator: {
          path: "./src/api/orval/fetchApiMutator.ts",
          name: "fetchApiMutator",
        },
      },
    },
  },
});
