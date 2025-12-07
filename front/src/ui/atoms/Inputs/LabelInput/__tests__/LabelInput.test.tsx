import { render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { LabelInput } from '../LabelInput';

describe('LabelInput', () => {
  it('renders a label with correct input link', () => {
    render(<LabelInput inputId="test-input">Label text</LabelInput>);
    const label = screen.getByText('Label text').closest('label');
    expect(label).toBeInTheDocument();
    expect(label).toHaveAttribute('for', 'test-input');
  });

  it('applies custom className in addition to default class', () => {
    render(
      <LabelInput inputId="test-input" className="custom-class">
        Label text
      </LabelInput>,
    );
    const label = screen.getByText('Label text').closest('label');
    expect(label).toHaveClass('ds-label-input custom-class');
  });

  it('shows required asterisk when required', () => {
    render(
      <LabelInput inputId="test-input" required>
        Label text
      </LabelInput>,
    );
    const asterisk = screen.getByText('*');
    expect(asterisk).toBeInTheDocument();
    expect(asterisk).toHaveAttribute('aria-hidden', 'true');
  });

  it('does not show asterisk when not required', () => {
    render(<LabelInput inputId="test-input">Label text</LabelInput>);
    expect(screen.queryByText('*')).toBeNull();
  });
});
