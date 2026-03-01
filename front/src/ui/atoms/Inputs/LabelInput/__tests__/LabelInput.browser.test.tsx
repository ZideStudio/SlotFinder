import { expect, test } from 'vitest';
import { render } from 'vitest-browser-react';
import { LabelInput } from '../LabelInput';

test('LabelInput visual snapshot (not required)', async () => {
  const { container } = await render(<LabelInput inputId="test-input">Label text</LabelInput>);

  await expect(container).toMatchScreenshot('label-input-not-required');
});

test('LabelInput visual snapshot (required)', async () => {
  const { container } = await render(
    <LabelInput inputId="test-input" required>
      Label text
    </LabelInput>,
  );

  await expect(container).toMatchScreenshot('label-input-required');
});
