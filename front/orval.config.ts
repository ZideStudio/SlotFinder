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
