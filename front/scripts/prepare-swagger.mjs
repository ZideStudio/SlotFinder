// oxlint-disable import/no-nodejs-modules
import { mkdirSync, readFileSync, writeFileSync } from "node:fs";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";

const directoryName = dirname(fileURLToPath(import.meta.url));

const swagger = JSON.parse(
  readFileSync(join(directoryName, "../../back/docs/swagger.json"), "utf-8"),
);

// Remove securityDefinitions to work around a crash in @scalar/openapi-parser
// (cookie-based apiKey is non-standard in Swagger 2.0 and triggers a TypeError
// In the validator's error-formatter). Security schemes are irrelevant for
// Client code generation.
delete swagger.securityDefinitions;

// Strip the /api base prefix from all paths so the fetch mutator can prepend
// FRONT_BACKEND_URL, which is fixed at build time.
const normalizedPaths = {};
for (const [path, value] of Object.entries(swagger.paths)) {
  const normalized = path.startsWith("/api/") ? path.slice(4) : path;
  normalizedPaths[normalized] = value;
}
swagger.paths = normalizedPaths;

const outputDir = join(directoryName, "../src/api/orval");
mkdirSync(outputDir, { recursive: true });

writeFileSync(
  join(outputDir, "swagger.json"),
  JSON.stringify(swagger, null, 2),
);

console.log("✅ Swagger prepared for Orval code generation");
