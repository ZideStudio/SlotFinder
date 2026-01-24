import { ColorInputAtom } from '../ColorInputAtom';
import { fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';

describe('ColorInput', () => {
  it('renders the color input with default value', () => {
    render(<ColorInputAtom name="color" />);

    expect(screen.getByText('#FF0000')).toBeInTheDocument();

    const input = screen.getByLabelText('#FF0000');
    expect(input).toBeInTheDocument();
    expect(input).toHaveValue('#ff0000');
  });

  it('updates the value when a new color is selected', () => {
    render(<ColorInputAtom name="color" />);

    const input = screen.getByLabelText('#FF0000');

    fireEvent.change(input, {
      target: { value: '#00ff00' },
    });

    expect(screen.getByText('#00ff00')).toBeInTheDocument();
    expect(input).toHaveValue('#00ff00');
  });

  it('applies extra props to the input', () => {
    render(<ColorInputAtom name="color" disabled />);

    const input = screen.getByLabelText('#FF0000');

    expect(input).toBeDisabled();
  });

  it('forwards required prop to the input', () => {
    render(<ColorInputAtom name="color" />);

    const input = screen.getByLabelText('#FF0000');

    expect(input).toHaveProperty('required');
  });

  it('forwards data attributes to the input', () => {
    render(<ColorInputAtom name="color" data-testid="color-input" />);

    const input = screen.getByTestId('color-input');

    expect(input).toBeInTheDocument();
  });

  it('forwards aria attributes to the input', () => {
    render(<ColorInputAtom name="color" aria-label="color picker" />);

    const input = screen.getByLabelText('color picker');

    expect(input).toBeInTheDocument();
  });
});
