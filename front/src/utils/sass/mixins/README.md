# Mixins

This folder contains reusable Sass mixins to apply consistent and complex styles across the project.

## Files

- **\_responsive.scss**: Mixins for handling media queries.
- **\_animations.scss**: Mixins for animations.

## Example

```scss
/* _responsive.scss */
@mixin responsive($breakpoint) {
  @if $breakpoint == small {
    @media (max-width: 576px) {
      @content;
    }
  } @else if $breakpoint == medium {
    @media (max-width: 768px) {
      @content;
    }
  }
}

/* _animations.scss */
@mixin fade-in {
  animation: fadeIn 1s ease-in-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}
```
