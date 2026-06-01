# 🎭 Mocks

The "mocks" folder contains mock HTTP request handlers and fixtures using [Mock Service Worker (MSW)](https://mswjs.io/). This allows the application to simulate API responses without relying on a real backend server.

## 📑 Table of Contents

- [Folder Organization](#folder-organization)
- [Handlers](#handlers)
- [Fixtures](#fixtures)
- [Learn More](#learn-more)

## <span id="folder-organization">Folder Organization</span>

- **`browser.ts`** : Sets up MSW for the browser environment (development & browser tests)
- **`server.ts`** : Sets up MSW for Node.js environment (unit tests)
- **`handlers/`** : Request handler definitions organized by domain
- **`fixtures/`** : Reusable mock data used by handlers

## <span id="handlers">Handlers</span>

Handlers define how to intercept and respond to specific HTTP requests.

### Naming Convention

Handlers follow the `{httpMethod}{ResourceName}{statusCode}` naming convention:

- **`httpMethod`** : HTTP verb in camelCase (`get`, `post`, `put`, `delete`, etc.)
- **`ResourceName`** : Resource name in PascalCase (e.g., `Account`, `AuthStatus`, `TokenRefresh`)
- **`statusCode`** : Expected HTTP status code (e.g., `200`, `201`, `400`, `401`)

For non-HTTP errors (e.g., network failures), replace the status code with a descriptive suffix:

```typescript
export const getAuthStatus200 = ...;            // GET /auth/status → 200
export const getAuthStatus401 = ...;            // GET /auth/status → 401
export const postTokenRefresh200 = ...;         // POST /auth/refresh → 200
export const postTokenRefreshNetworkError = ...; // POST /auth/refresh → network error
```

### Creating a Handler

```typescript
// mocks/handlers/authStatusHandlers.ts
import {
  getAuthStatus200Fixture,
  getAuthStatus401Fixture,
} from "@Mocks/fixtures/authStatusFixtures";
import { delay, http, HttpResponse } from "msw";

export const getAuthStatus200 = http.get(
  `${import.meta.env.FRONT_BACKEND_URL}/v1/auth/status`,
  async () => {
    await delay();
    return HttpResponse.json(getAuthStatus200Fixture, { status: 200 });
  },
);

export const getAuthStatus401 = http.get(
  `${import.meta.env.FRONT_BACKEND_URL}/v1/auth/status`,
  async () => {
    await delay();
    return HttpResponse.json(getAuthStatus401Fixture, { status: 401 });
  },
);
```

### Registering Handlers

Add handlers to `browser.ts` and `server.ts`:

```typescript
// browser.ts
import { setupWorker } from "msw/browser";

export const worker = setupWorker(getAuthStatus200, getAuthStatus401);
```

## <span id="fixtures">Fixtures</span>

Fixtures provide reusable mock data for handlers. **Always use the same types from your application** to ensure type consistency and catch API changes automatically.

### Naming Convention

Fixtures follow the `{httpMethod}{ResourceName}{statusCode}Fixture` naming convention:

- **`httpMethod`** : HTTP verb in camelCase (`get`, `post`, `put`, `delete`, etc.)
- **`ResourceName`** : Resource name in PascalCase (e.g., `Account`, `AuthStatus`)
- **`statusCode`** : Expected HTTP status code (e.g., `200`, `201`, `400`, `401`)

```typescript
export const getAccount200Fixture = ...;   // GET /account → 200
export const postAccount201Fixture = ...;  // POST /account → 201
export const postAccount400Fixture = ...;  // POST /account → 400
```

### Creating a Fixture

```typescript
// mocks/fixtures/accountFixtures.ts
import type { ErrorResponseType } from "@Front/types/api.types";
import type {
  SignUpResponseType,
  SignUpErrorCodeType,
} from "@Front/types/Authentication/signUp/signUp.types";

export const postAccount201Fixture: SignUpResponseType = {
  access_token: "1234567890abcdef",
  email: "test@example.com",
  id: "123456",
  userName: "test_user",
};

export const postAccount400Fixture: ErrorResponseType<SignUpErrorCodeType> = {
  code: "USERNAME_ALREADY_TAKEN",
};
```

When an application type changes, TypeScript will flag the fixture as invalid, reminding you to update it accordingly.

### Using Fixtures in Handlers

```typescript
// mocks/handlers/accountHandlers.ts
import {
  postAccount201Fixture,
  postAccount400Fixture,
} from "../fixtures/accountFixtures";

export const postAccount201 = http.post("/api/v1/account", async () => {
  return HttpResponse.json(postAccount201Fixture, { status: 201 });
});

export const postAccount400 = http.post("/api/v1/account", async () => {
  return HttpResponse.json(postAccount400Fixture, { status: 400 });
});
```

## <span id="learn-more">Learn More</span>

- [Mock Service Worker Documentation](https://mswjs.io/docs)
- [HTTP Handlers](https://mswjs.io/docs/api/http)
