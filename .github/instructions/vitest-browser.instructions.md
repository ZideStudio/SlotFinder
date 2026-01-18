---
applyTo: '**/*.browser.test.{ts,tsx,js,jsx}'
description: Instructions for writing visual regression and browser integration tests with Vitest Browser Mode
---

# Vitest Browser Testing Instructions (Front)

## Purpose

Ensure quality, clarity, and maintainability of browser tests in the front-end using Vitest Browser Mode with Playwright. These tests run components in a real Chromium browser environment, enabling visual regression testing and real browser interaction testing.

## When to Use Browser Tests

**✅ Use browser tests for:**
- **Visual regression testing** with screenshot snapshots (atoms, molecules, organisms)
- Components with real browser API dependencies
- Keyboard navigation and focus management testing
- Accessibility testing in real browser environment
- Complex user interactions and animations

**Specifically for Atoms & Molecules:**
- Use browser tests for **visual regression testing** when visual consistency is critical
- Examples: Button states, Badge variants, Form inputs, Labels, Icons
- These components are reused everywhere, so visual regressions should be caught early
- Browser tests ensure styles and layout render correctly across the real browser

**❌ Don't use browser tests for:**
- Pure logic testing (use unit tests instead)
- Tests not requiring browser environment

## Component Level Testing Strategy

| Component Level | Unit Tests (`*.test.tsx`) | Browser Tests (`*.browser.test.tsx`) |
|-----------------|---------------------------|-------------------------------------|
| **Atoms** | Logic & rendering | ✅ Visual regression (recommended) |
| **Molecules** | Integration & logic | ✅ Visual regression & interactions |
| **Organisms** | Complex logic | ✅ Visual regression & E2E-like tests |
| **Pages** | Full integration | ❌ Use E2E (Playwright) instead |

## File Naming and Organization

- Browser test files must use the suffix: `ComponentName.browser.test.tsx`
- Place browser tests in the `__tests__` folder next to unit tests
- Example structure:
  ```
  📁 Button
  ├── 📁 __tests__
  │   ├── Button.test.tsx          # Unit tests
  │   ├── Button.browser.test.tsx  # Browser/visual tests
  ├── Button.tsx
  ├── Button.module.scss
  ```

## Structure and Organization

- Each browser test file should focus on one component
- Group related visual tests using multiple `test()` calls
- Test different component states separately (default, disabled, error, etc.)
- Each test should be independent and not rely on other tests

## Imports and Tools

```typescript
import { expect, test } from 'vitest';
import { render } from 'vitest-browser-react';
```

## Visual Regression Testing

### Screenshot Capture Best Practices

- Always disable animations: `animations: 'disabled'`
- Keep full page captures off unless necessary: `fullPage: false`
- Verify component is rendered before capturing screenshot

### Snapshot Comparison

- Use strict threshold: `threshold: 0` to detect any pixel-level change
- Name snapshots explicitly and clearly: `component-state.png`
- Store snapshots alongside test file in a `__snapshots__` directory

### Example: Visual Regression Test for Atom

```typescript
// src/ui/atoms/Button/__tests__/Button.browser.test.tsx
import { expect, test } from 'vitest';
import { render } from 'vitest-browser-react';
import { Button } from '../Button';

// Test all visual states of the Button atom
test('Button visual snapshot (default state)', async () => {
  const { getByRole } = await render(<Button>Click me</Button>);
  
  // Verify component is rendered correctly
  await expect.element(getByRole('button')).toBeInTheDocument();
  
  // Capture and compare screenshot
  await expect(getByRole('button')).toMatchScreenshot('button-default');
});

test('Button visual snapshot (disabled state)', async () => {
  const { getByRole } = await render(<Button disabled>Disabled</Button>);
  
  await expect.element(getByRole('button')).toBeDisabled();
  
  // Capture and compare screenshot
  await expect(getByRole('button')).toMatchScreenshot('button-disabled');
});

test('Button visual snapshot (loading state)', async () => {
  const { getByRole } = await render(<Button isLoading>Loading...</Button>);
  
  await expect.element(getByRole('button')).toBeDisabled();
  
  // Capture and compare screenshot
  await expect(getByRole('button')).toMatchScreenshot('button-loading');
});
```

### Example: Visual Regression Test for Molecule

```typescript
// src/ui/molecules/FormField/__tests__/FormField.browser.test.tsx
import { expect, test } from 'vitest';
import { render } from 'vitest-browser-react';
import { FormField } from '../FormField';

test('FormField visual snapshot (default)', async () => {
  const { container } = await render(
    <FormField label="Email" required>
      <input type="email" />
    </FormField>
  );
  
  // Capture and compare screenshot
  await expect(container).toMatchScreenshot('form-field-default');
});

test('FormField visual snapshot (error state)', async () => {
  const { container, getByRole } = await render(
    <FormField label="Email" error="Invalid email">
      <input type="email" />
    </FormField>
  );
  
  await expect.element(getByRole('alert')).toHaveAccessibleName(/Invalid email/i);
  
  // Capture and compare screenshot
  await expect(container).toMatchScreenshot('form-field-error');
});
```

## Accessibility Testing in Browser

- Use accessible selectors with `vitest-browser-react` queries
- Test keyboard navigation and focus management
- Verify ARIA attributes and roles in real browser
- Test with real screen reader behavior expectations

### Accessibility Example

```typescript
test('Form field keyboard navigation', async () => {
  const { getByLabelText, getByRole } = await render(
    <FormField label="Name" required>
      <input type="text" />
    </FormField>
  );
  
  // Test keyboard accessible name
  await expect.element(getByLabelText('Name')).toHaveAccessibleName('Name');
  
  // Test required attribute visibility
  const requiredIndicator = getByRole('img', { hidden: true });
  await expect.element(requiredIndicator).toHaveAttribute('aria-hidden', 'true');
});
```

## Snapshot Management Workflow

### First Time - Create Baseline Snapshots

```bash
npm run test:browser:update
```

This generates the initial visual baseline. Always review generated screenshots to ensure they are correct.

### Development - Run Tests Against Baseline

```bash
npm run test:browser
```

Tests will fail if rendered output differs from baseline snapshots.

### After Intentional Changes - Update Snapshots

```bash
npm run test:browser:update
```

Always visually validate changes before committing updated snapshots.

### Watch Mode for Development

```bash
npm run test:browser:watch
```

Useful for interactive development and immediate feedback.

## Best Practices

- **Isolate components**: Render only the component being tested, not entire pages
- **Keep tests focused**: Each test should verify one specific visual state or behavior
- **Use consistent naming**: Name snapshots to reflect the tested state (e.g., `component-state.png`)
- **Document visual importance**: Add comments explaining why visual testing is critical for the component
- **Review snapshots in version control**: Always review snapshot changes in PRs
- **Avoid flaky tests**: Disable animations, set consistent viewport, avoid timing-dependent assertions
- **Test meaningful variations**: Focus on states that affect visual appearance (disabled, error, loading, etc.)
- **For atoms & molecules**: Test all visual variants and states (default, hover, disabled, error, loading, etc.)
- **For organisms**: Combine visual tests with interaction tests for more realistic scenarios

## Example: Complete Browser Test

```typescript
import { expect, test } from 'vitest';
import { render } from 'vitest-browser-react';
import { page } from 'vitest/browser';
import { Card } from '../Card';

describe('Card', () => {
  test('Card visual snapshot (default)', async () => {
    const { getByRole } = await render(
      <Card title="Test Card">
        Content goes here
      </Card>
    );
    
    await expect.element(getByRole('heading')).toHaveAccessibleName('Test Card');
    
    const screenshot = await page.screenshot({
      fullPage: false,
      animations: 'disabled',
    });
    
    expect(screenshot).toMatchSnapshot({
      name: 'card-default.png',
      threshold: 0,
    });
  });

  test('Card visual snapshot (with loading state)', async () => {
    const { getByRole } = await render(
      <Card title="Loading" isLoading>
        Content
      </Card>
    );
    
    const loadingIndicator = getByRole('status');
    await expect.element(loadingIndicator).toHaveAccessibleName(/loading/i);
    
    const screenshot = await page.screenshot({
      fullPage: false,
      animations: 'disabled',
    });
    
    expect(screenshot).toMatchSnapshot({
      name: 'card-loading.png',
      threshold: 0,
    });
  });

  test('Card visual snapshot (with error state)', async () => {
    const { getByRole } = await render(
      <Card title="Error" error="Something went wrong">
        Content
      </Card>
    );
    
    const alert = getByRole('alert');
    await expect.element(alert).toHaveAccessibleName(/error/i);
    
    const screenshot = await page.screenshot({
      fullPage: false,
      animations: 'disabled',
    });
    
    expect(screenshot).toMatchSnapshot({
      name: 'card-error.png',
      threshold: 0,
    });
  });
});
```

## Common Issues and Solutions

### Snapshots Regenerating Unexpectedly
- Check that `UPDATE_SNAPSHOTS` environment variable is not set when running `test:browser`
- Ensure you're using `test:browser` (not `test:browser:update`) for normal test runs

### Flaky Visual Tests
- Disable animations with `animations: 'disabled'`
- Avoid assertions that depend on timing
- Use explicit assertions before screenshot to ensure rendered state

### Style Not Applied in Tests
- Verify CSS/SCSS imports in component are correct
- Check that styles are loaded in Vitest browser configuration
- Ensure CSS modules are properly imported

## References

- [Vitest Browser Mode Documentation](https://vitest.dev/guide/browser.html)
- [Playwright Documentation](https://playwright.dev/)
- [Testing Library Documentation](https://testing-library.com/)
- [Accessibility instructions](./a11y.instructions.md)
- [Vitest classic instructions](./vitest.instructions.md)
