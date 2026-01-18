import { ColorInput } from '../ColorInput';
import { fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';

describe('ColorInput', () => {
  it('renders the color input with default value', () => {
    render(<ColorInput name="color" />);

    expect(screen.getByText('#FF0000')).toBeInTheDocument();

    const input = screen.getByLabelText('#FF0000') as HTMLInputElement;
    expect(input).toBeInTheDocument();
    expect(input.value).toBe('#ff0000');
  });

  it('updates the value when a new color is selected', () => {
    render(<ColorInput name="color" />);

    const input = screen.getByLabelText('#FF0000') as HTMLInputElement;

    fireEvent.change(input, {
      target: { value: '#00FF00' },
    });

    expect(screen.getByText('#00ff00')).toBeInTheDocument();
    expect(input.value).toBe('#00ff00');
  });
});
