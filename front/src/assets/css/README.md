# Global CSS Assets

This folder contains the global style foundations used by the front-end application.

It defines:

- design tokens (colors and theme variables)
- spacing and dimension tokens generated from a Sass scale
- the native grid/subgrid layout system
- the global reset and baseline CSS
- shared input base styling for UI atoms

## Files Overview

- `index.ts`: single entry point imported by the app and Storybook
- `tokens.css`: global color and state custom properties
- `_dimension.scss`: generated margin, padding, and gap custom properties
- `_grid.scss`: grid and subgrid layout utilities
- `reset.css`: element reset and typography defaults
- `global.css`: minimal root-level global variables
- `_inputs.scss`: reusable base input style placeholder

## Load Order and Cascade

The global entry point is `index.ts` and it imports files in this order:

1. `_dimension.scss`
2. `_grid.scss`
3. `global.css`
4. `reset.css`
5. `tokens.css`

This order means:

- layout tokens and utilities are available early
- reset rules normalize defaults before most component-specific styles are applied
- token declarations are available globally after load

## Design Tokens (`tokens.css`)

`tokens.css` declares global custom properties on `:root`.

### Color families

- Primary: `--color-primary` and related light/dark/soft variants
- Accent: `--color-accent` and related light/dark/soft variants
- Neutral greys: `--color-grey-100`, `--color-grey-200`, `--color-grey-300`, `--color-grey-400`, `--color-grey-900`
- Base colors: `--color-white`, `--color-black`, `--color-blue`

### State colors

- `--color-focus`
- `--color-error`
- `--color-success`
- `--color-warning`

Notes:

- some derived colors use modern CSS color functions from HSL source tokens
- all tokens are globally available and can be consumed from plain CSS or SCSS

## Dimension System (`_dimension.scss`)

The dimension system generates spacing variables from an 8px step converted to rem.

Base conversion:

- 1 unit = 8px = 0.5rem (root font size 16px)

Generated token groups:

- margins: `--margin-0` through `--margin-15`
- paddings: `--padding-0` through `--padding-4`
- gaps: `--gap-0` through `--gap-4`

Half-step tokens are also generated for selected values:

- `--margin-0-5`, `--margin-2-5`, `--margin-3-5`
- `--padding-0-5`, `--padding-2-5`, `--padding-3-5`
- `--gap-0-5`, `--gap-2-5`, `--gap-3-5`

Implementation detail:

- token values are generated via Sass mixins and helper functions
- this keeps spacing consistent across components and layout utilities

## Native Grid System (`_grid.scss`)

This project uses a native CSS Grid system powered by utility classes and CSS custom properties.

Goal:

- keep layout logic in stylesheets, not in component props

### Core classes

- `.grid`: root layout container with responsive column setup
- `.subgrid`: nested container aligned with parent columns

### Item positioning

Direct children of `.grid` and `.subgrid` are positioned with:

- `--cols`: number of columns to span
- `--start`: grid-column start position

Applied rule:

- `grid-column: var(--start) / span var(--cols)`

### Spacing variables used by grid

- `--gap` maps to column gap
- `--row-gap` maps to row gap
- `--margin-h` maps to horizontal container margin

### Mobile-first defaults

- `--cols: 4`
- `--gap: var(--gap-2)`
- `--row-gap: var(--gap-2)`
- `--margin-h: var(--margin-2)`

### Breakpoint behavior

Responsive values are applied through Sass mixins:

- tablet: 8 columns, horizontal margin token 4
- desktop-small: 12 columns, larger horizontal margin
- desktop-medium: 12 columns, larger horizontal margin
- desktop-large: 12 columns, centered container (`--margin-h: auto`)

The grid container max-width is 1440px (converted to rem).

### Subgrid fallback

When native subgrid is supported:

- `grid-template-columns: subgrid`

When not supported:

- fallback to `repeat(var(--cols), minmax(0, 1fr))`
- use `column-gap: var(--gap)`

### Grid examples

#### 1) Basic page layout with `.grid`

```html
<main class="grid page-layout">
  <section class="hero">Hero</section>
  <section class="content">Content</section>
  <aside class="sidebar">Sidebar</aside>
</main>
```

```scss
.page-layout {
  .hero {
    --cols: 4;
    --start: 1;
  }

  .content {
    --cols: 4;
    --start: 1;
  }

  .sidebar {
    --cols: 4;
    --start: 1;
  }
}
```

#### 2) Nested alignment with `.subgrid`

```html
<section class="grid profile-page">
  <article class="subgrid profile-card">
    <h2 class="profile-card__title">Profile</h2>
    <p class="profile-card__summary">Short summary</p>
    <div class="profile-card__actions">Actions</div>
  </article>
</section>
```

```scss
.profile-card {
  --cols: 4;
  --start: 1;

  &__title,
  &__summary,
  &__actions {
    --cols: 4;
    --start: 1;
  }
}
```

#### 3) Responsive placement with SCSS mixins

```tsx
<main className="grid auth-page">
  <section className="auth-page__main subgrid">
    <h1 className="auth-page__title">Sign in</h1>
    <div className="auth-page__form">{/* form */}</div>
    <nav className="auth-page__oauth">{/* providers */}</nav>
  </section>
</main>
```

```scss
@use "../../utils/sass/mixins/responsive";

.auth-page {
  &__main {
    row-gap: var(--row-gap);
  }

  &__form {
    --cols: 4;
    --start: 1;

    @include responsive.mediaExceedsWidth("tablet") {
      --cols: 2;
      --start: 4;
    }

    @include responsive.mediaExceedsWidth("desktop-small") {
      --start: 6;
    }
  }

  &__oauth {
    --cols: 4;
    --start: 1;

    @include responsive.mediaExceedsWidth("tablet") {
      --start: 3;
    }

    @include responsive.mediaExceedsWidth("desktop-small") {
      --start: 5;
    }
  }
}
```

## Reset and Baseline (`reset.css` + `global.css`)

### `reset.css`

- strips default margin/padding/border from common elements
- normalizes font inheritance and smoothing
- sets HTML5 semantic elements to block display for legacy consistency
- removes list bullets and quote pseudo-content defaults
- normalizes table collapse/spacing behavior
- applies a monospace stack to `code`

### `global.css`

- currently defines `--font-size-base: 16` at root level

## Shared Input Base (`_inputs.scss`)

`_inputs.scss` defines a Sass placeholder selector:

- `%inputBase`

It is consumed in input atoms with Sass `@extend`, for example:

- `TextInputAtom.scss` extends `%inputBase`

### What `%inputBase` provides

- full-width input layout
- base shadow, radius, and typography
- responsive font size increase on desktop-small and above
- placeholder color handling
- focus, hover, and filled state visual behavior
- error visuals through `[aria-invalid="true"]`

This ensures a single source of truth for text-like input controls.

## Practical Usage Guidelines

1. Import styles once via `@Front/assets/css`.
2. Use token variables instead of hardcoded values whenever possible.
3. Use `.grid` for page-level layout and `.subgrid` for nested alignment.
4. Set placement through `--cols` and `--start` on item classes.
5. Reuse `%inputBase` for any new text-like input atom.
6. Keep spacing consistent with `--margin-*`, `--padding-*`, and `--gap-*` tokens.

## Related Files

- Grid implementation: `front/src/assets/css/_grid.scss`
- Token palette: `front/src/assets/css/tokens.css`
- Input base style: `front/src/assets/css/_inputs.scss`
- Global importer: `front/src/assets/css/index.ts`
