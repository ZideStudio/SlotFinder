# Route Authentication Registry

## Overview

The `routeAuthRegistry` provides a centralized way to track which routes require authentication. This approach avoids hardcoded path lists and integrates directly with route definitions.

## How It Works

### Registration

Routes that don't require authentication register themselves when their module is imported:

```typescript
// In route definition file (e.g., SignUp/routes.tsx)
import { routeAuthRegistry } from '@Front/utils/routeAuthRegistry';

// Register this route as not requiring authentication
routeAuthRegistry.registerUnauthenticatedPath('sign-up');

export const signUpRoutes: RouteObject = {
  path: 'sign-up',
  element: <SignUp />,
  handle: {
    mustBeAuthenticate: false,  // This should match the registration
  },
};
```

### Usage

The `isAuthenticatedPage()` utility uses the registry to determine if the current page requires authentication:

```typescript
import { isAuthenticatedPage } from '@Front/utils/isAuthenticatedPage';

// In fetchApi or other code
if (isAuthenticatedPage()) {
  // Current page requires authentication
  // Safe to trigger token refresh on 401
}
```

## Convention

- **Default:** All routes require authentication unless explicitly registered as unauthenticated
- **Unauthenticated routes:** Must call `routeAuthRegistry.registerUnauthenticatedPath()` AND set `handle.mustBeAuthenticate = false`
- **Authenticated routes:** Do nothing - authentication is required by default

## Benefits

1. **No Hardcoded Lists:** Routes register themselves, avoiding a centralized list that must be maintained
2. **Collocated Configuration:** Authentication requirements are defined alongside route configuration
3. **Type-Safe:** TypeScript ensures proper types throughout
4. **Test-Friendly:** No circular dependencies or route object imports in test environments
5. **Maintainable:** Adding new routes doesn't require changes to multiple files

## Adding New Routes

### Authenticated Route (Default)

Just create your route - no special action needed:

```typescript
export const dashboardRoutes: RouteObject = {
  path: 'dashboard',
  element: <Dashboard />,
  // Authentication required by default
};
```

### Unauthenticated Route

Register the path and set the handle property:

```typescript
import { routeAuthRegistry } from '@Front/utils/routeAuthRegistry';

routeAuthRegistry.registerUnauthenticatedPath('login');

export const loginRoutes: RouteObject = {
  path: 'login',
  element: <Login />,
  handle: {
    mustBeAuthenticate: false,
  },
};
```

## Current Registered Unauthenticated Routes

- `sign-up` - User registration page
- `oauth/callback` - OAuth callback handler

These match routes with `handle.mustBeAuthenticate = false` in their route configuration.
