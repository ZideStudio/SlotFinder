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

### Creating a Handler

```typescript
// mocks/handlers/userHandlers.ts
import { http, HttpResponse } from "msw";

export const getUserHandler = http.get("/api/users/:id", ({ params }) => {
  return HttpResponse.json({ id: params.id, name: "John Doe" });
});

export const createUserHandler = http.post(
  "/api/users",
  async ({ request }) => {
    const body = await request.json();
    return HttpResponse.json({ id: 1, ...body }, { status: 201 });
  },
);
```

### Registering Handlers

Add handlers to `browser.ts` and `server.ts`:

```typescript
// browser.ts
import { setupWorker } from "msw/browser";

export const worker = setupWorker(getUserHandler, createUserHandler);
```

## <span id="fixtures">Fixtures</span>

Fixtures provide reusable mock data for handlers. **Always use the same types from your application** to ensure type consistency and catch API changes automatically.

### Creating a Fixture

```typescript
// mocks/fixtures/userFixtures.ts
import { User } from "@Front/types/user"; // Use application types

export const userFixture: User = {
  id: 1,
  name: "John Doe",
  email: "john@example.com",
};
```

When the `User` type changes in your application, TypeScript will flag the fixture as invalid, reminding you to update it accordingly.

### Using Fixtures in Handlers

```typescript
// mocks/handlers/userHandlers.ts
import { userFixture } from "../fixtures/userFixtures";

export const getUserHandler = http.get("/api/users/1", () => {
  return HttpResponse.json(userFixture);
});
```

## <span id="learn-more">Learn More</span>

- [Mock Service Worker Documentation](https://mswjs.io/docs)
- [HTTP Handlers](https://mswjs.io/docs/api/http)
