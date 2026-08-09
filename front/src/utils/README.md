# 📁 utils

The utils folder contains shared front-end helpers and support utilities used across the application. It centralizes small, reusable pieces of logic so components and pages can stay focused on business behavior instead of repeating boilerplate.

## 📑 Table of Contents

- [Structure](#structure)
- [Available utilities](#available-utilities)
- [Best practices](#best-practices)

## 🧱 Structure

The folder is organized into four main areas:

- `constants/`: shared constants such as terminology and duration units
- `helpers/`: reusable pure functions for common UI and URL logic
- `sass/`: SCSS mixins and helpers for layout and spacing
- `testsUtils/`: test helpers for rendering components with providers, routing, and query clients

## 🛠️ Available utilities

### Helpers

#### `getClassName`

Builds a class name from a base class and optional modifiers.

```ts
import { getClassName } from "@Front/utils/helpers/getClassName";

const className = getClassName({
  defaultClassName: "button",
  modifiers: ["primary", "large"],
  className: "custom-button",
});
```

#### `getContrastTextColor`

Returns a readable text color (`#000000` or `#FFFFFF`) based on a background hex color.

```ts
import { getContrastTextColor } from "@Front/utils/helpers/getContrastTextColor";

const textColor = getContrastTextColor("#336699");
```

#### `isInternalUrl`

Validates that a URL is a safe internal route and blocks common redirect bypass patterns.

```ts
import { isInternalUrl } from "@Front/utils/helpers/isInternalUrl";

if (isInternalUrl("/dashboard")) {
  // safe internal navigation
}
```

### Constants

#### `terms.ts`

Exports the current terms version constant.

```ts
import { TERMS_VERSION } from "@Front/utils/constants/terms";
```

#### `units.ts`

Defines the supported duration units and their allowed ranges.

```ts
import { UNITS, UNIT_LIMITS } from "@Front/utils/constants/units";
```

### Test utilities

The test utilities provide shared rendering helpers for React Testing Library, React Query, router context, and application providers.

```ts
import {
  renderWithQueryClient,
  renderRoute,
} from "@Front/utils/testsUtils/customRender/customRender";
```

These helpers are especially useful for component and page tests that need authentication, loaders, toasts, or router state.

### Sass utilities

The Sass helpers under `sass/mixins/` provide reusable mixins for spacing, responsive behavior, and common layout patterns.

Typical files include:

- `_spacing.scss`
- `_responsive.scss`
- `_helpers.scss`

## ✅ Best practices

- Keep utilities small, focused, and side-effect free.
- Prefer pure functions with explicit return types when possible.
- Reuse existing helpers before introducing a new one.
- Keep imports consistent and use the most specific path available.
- Add or update tests alongside utility changes, especially for helpers in `helpers/__tests__/`.
- Avoid putting business logic or component-specific behavior inside shared utils.

This folder should remain a lightweight collection of reusable primitives that improve consistency across the codebase.
