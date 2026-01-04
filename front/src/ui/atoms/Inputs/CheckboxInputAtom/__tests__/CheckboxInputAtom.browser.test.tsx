import { expect, test } from 'vitest';
import { render } from 'vitest-browser-react';
import { CheckboxInputAtom } from '../CheckboxInputAtom';

test('CheckboxInputAtom visual snapshot', async () => {
  const { container } = await render(<CheckboxInputAtom name="test-input" />);

  await expect(container).toMatchScreenshot('checkbox-input-atom-default');
});

test('CheckboxInputAtom visual snapshot (checked)', async () => {
  const { container } = await render(<CheckboxInputAtom name="test-input" checked />);

  await expect(container).toMatchScreenshot('checkbox-input-atom-checked');
});

test('CheckboxInputAtom visual snapshot (disabled)', async () => {
  const { container } = await render(<CheckboxInputAtom name="test-input" disabled />);

  await expect(container).toMatchScreenshot('checkbox-input-atom-disabled');
});

test('CheckboxInputAtom visual snapshot (checked and disabled)', async () => {
  const { container } = await render(<CheckboxInputAtom name="test-input" checked disabled />);

  await expect(container).toMatchScreenshot('checkbox-input-atom-checked-disabled');
});

test('CheckboxInputAtom visual snapshot (unchecked and on error)', async () => {
  const { container } = await render(<CheckboxInputAtom name="test-input" aria-invalid="true" />);

  await expect(container).toMatchScreenshot('checkbox-input-atom-error');
});

test('CheckboxInputAtom visual snapshot (checked and on error)', async () => {
  const { container } = await render(<CheckboxInputAtom name="test-input" checked aria-invalid="true" />);

  await expect(container).toMatchScreenshot('checkbox-input-atom-checked-error');
});
