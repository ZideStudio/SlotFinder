# Native Grid System

This project uses a native CSS Grid system powered by utility classes and CSS custom properties (tokens).

The goal is to keep layout logic in stylesheets, not in React props.

## Where It Is Defined

- Grid styles: `front/src/components/Grid/_grid.scss`
- Global import: `front/src/assets/css/index.ts`

## Core Concepts

### 1) `grid` for root layout containers

Apply the `grid` class to a container that defines the main column structure.

```html
<main class="grid">
  <section class="hero"></section>
  <section class="content"></section>
  <aside class="sidebar"></aside>
</main>
```

Each direct child can be positioned with:

- `--cols`: number of columns to span
- `--start`: column where the item starts

### 2) `subgrid` for nested layout containers

Apply `subgrid` when a nested container must inherit the parent grid columns.

```html
<section class="grid">
  <div class="panel subgrid">
    <div class="panel__title"></div>
    <div class="panel__body"></div>
  </div>
</section>
```

Children inside `.subgrid` use the same `--cols` and `--start` logic.

## Available Tokens

### Item positioning tokens

Set these tokens on grid items (direct children of `.grid` or `.subgrid`):

- `--cols` (default behavior depends on context)
- `--start` (default: `auto`)

Example:

```scss
.page-content {
  --cols: 4;
  --start: 1;
}
```

### Spacing tokens

The system supports independent horizontal and vertical spacing:

- `--gap`: column gap token (mapped to `column-gap`)
- `--row-gap`: row gap token (mapped to `row-gap`)
- `--margin-h`: horizontal container margin

Example override on a specific grid:

```scss
.custom-layout {
  --gap: var(--gap-2);
  --row-gap: var(--gap-1);
}
```

## Mobile-First Behavior

The grid is mobile-first.

Default values on mobile:

- `--cols: 4`
- `--gap: var(--gap-2)`
- `--row-gap: var(--gap-2)`
- `--margin-h: var(--margin-2)`

At larger breakpoints, values are updated automatically by the grid stylesheet:

- `mobile` (`>= 0px`): 4 columns
- `tablet` (`>= 668px`): 8 columns
- `desktop-small` (`>= 1024px`): 12 columns
- `desktop-medium` (`>= 1280px`): 12 columns
- `desktop-large` (`>= 1600px`): 12 columns

Use project mixins for responsive overrides:

```scss
@use "../../utils/sass/mixins/responsive";

.card {
  --cols: 4;

  @include responsive.media-exceeds-width("tablet") {
    --cols: 2;
    --start: 4;
  }

  @include responsive.media-exceeds-width("desktop-small") {
    --start: 6;
  }
}
```

## Recommended Usage Pattern

1. Put `grid` on the page/container root.
2. Put `subgrid` on nested containers that must align with parent columns.
3. Control placement with `--cols` and `--start` on each item class.
4. Control vertical rhythm with `--row-gap` when needed.
5. Prefer SCSS responsive mixins over hardcoded media query values.

## Full Example

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

    @include responsive.media-exceeds-width("tablet") {
      --cols: 2;
      --start: 4;
    }

    @include responsive.media-exceeds-width("desktop-small") {
      --start: 6;
    }
  }

  &__oauth {
    --cols: 4;
    --start: 1;

    @include responsive.media-exceeds-width("tablet") {
      --start: 3;
    }

    @include responsive.media-exceeds-width("desktop-small") {
      --start: 5;
    }
  }
}
```
