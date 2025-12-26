import { render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { InputErrorMessage } from '../InputErrorMessage';

describe('', () => {
  it('renders an error message with correct message', () => {
    render(<InputErrorMessage message="test-input"></InputErrorMessage>);
    const errorMessage = screen.getByText('test-input');
    expect(errorMessage).toBeInTheDocument();
    expect(errorMessage).toHaveTextContent('test-input');
  });
  it('applies custom id when provided', () => {
    render(<InputErrorMessage message="test-input" id="test-id"></InputErrorMessage>);
    const errorMessage = screen.getByText('test-input');
    expect(errorMessage).toHaveAttribute('id', 'test-id');
  });
  

});
