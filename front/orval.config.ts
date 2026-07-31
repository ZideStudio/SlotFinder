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
        runtime: "import.meta.env.FRONT_BACKEND_URL",
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
