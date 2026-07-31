# 📁 api

This folder contains the frontend API layer generated with Orval.

## Goal

- Generate typed API clients directly from backend Swagger
- Keep frontend API calls aligned with the backend contract
- Centralize request behavior (auth refresh, errors, headers) through a single mutator

## Structure

- `generated/`: Orval output (endpoints + schemas). Do not edit manually.
- `orval/swagger.json`: prepared Swagger input used for generation.
- `orval/fetchApiMutator.ts`: shared custom fetch logic used by generated clients.
- `tokenRefreshManager.ts`: token refresh logic used by the mutator.

## Generate API clients

Run:

```bash
npm run generate:api
```

The command:

1. Reads backend Swagger from `../back/docs/swagger.json`
2. Prepares it via `scripts/prepare-swagger.mjs`
3. Generates clients into `src/api/generated` using Orval

## Use generated functions

Import functions and types from generated files instead of hand-written API wrappers.

```typescript
import { postV1Account } from "@Front/api/generated/account/account";
```

## Best practices

- Re-run `npm run generate:api` after backend API changes.
- Do not manually edit files under `src/api/generated`.
- Keep custom request behavior in `orval/fetchApiMutator.ts`.
