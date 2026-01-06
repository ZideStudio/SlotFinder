import { render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { InputErrorMessage } from '../InputErrorMessage';

describe('InputErrorMessage', () => {
  it('renders an error message with correct message', () => {
    render(<InputErrorMessage>test-input</InputErrorMessage>);
    const errorMessage = screen.getByText('test-input');
    expect(errorMessage).toBeInTheDocument();
  });

  it('applies custom id when provided', () => {
    render(<InputErrorMessage id="test-id">test-input</InputErrorMessage>);
    const errorMessage = screen.getByText('test-input');
    expect(errorMessage).toHaveAttribute('id', 'test-id');
  });

  it('does not render anything when children is empty', () => {
    const { container } = render(<InputErrorMessage>{null}</InputErrorMessage>);
    expect(container.firstChild).toBeNull();
  });
});
