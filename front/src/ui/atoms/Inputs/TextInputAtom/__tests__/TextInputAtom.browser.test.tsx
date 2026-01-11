import { expect, test } from 'vitest';
import { render } from 'vitest-browser-react';
import { page } from 'vitest/browser';
import { TextInputAtom } from '../TextInputAtom';

test('TextInputAtom visual snapshot', async () => {
  const { container } = await render(<TextInputAtom name="test-input" placeholder="Enter text here..." />);

  await expect(container).toMatchScreenshot('text-input-atom-default');
});

test('TextInputAtom visual snapshot with value', async () => {
  const { container } = await render(<TextInputAtom name="test-input" defaultValue="Sample text" />);

  await expect(container).toMatchScreenshot('text-input-atom-with-value');
});

test('TextInputAtom visual snapshot on focus', async () => {
  const { container } = await render(<TextInputAtom name="test-input" />);

  const input = page.getByRole('textbox');
  await input.click();

  await expect(container).toMatchScreenshot('text-input-atom-on-focus');
});

test('TextInputAtom visual snapshot on hover', async () => {
  const { container } = await render(<TextInputAtom name="test-input" />);

  const input = page.getByRole('textbox');
  await input.hover();

  await expect(container).toMatchScreenshot('text-input-atom-on-hover');
});
