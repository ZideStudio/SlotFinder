import { ColorInputAtom } from '../ColorInputAtom';
import { fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';

describe('ColorInput', () => {
  it('renders the color input with default value', () => {
    render(<ColorInputAtom name="color" description="Choisir une couleur" />);

    expect(screen.getByText('Choisir une couleur')).toBeInTheDocument();

    const input = screen.getByLabelText('Choisir une couleur');
    expect(input).toBeInTheDocument();
    expect(input).toHaveValue('#ffffff');
  });

  it('updates the value when a new color is selected', () => {
    render(<ColorInputAtom name="color" description="Choisir une couleur" />);

    const input = screen.getByLabelText('Choisir une couleur');

    fireEvent.change(input, {
      target: { value: '#00ff00' },
    });

    expect(input).toHaveValue('#00ff00');
  });

  it('applies extra props to the input', () => {
    render(<ColorInputAtom name="color" disabled description="La couleur" />);

    const input = screen.getByLabelText('La couleur');

    expect(input).toBeDisabled();
  });

  it('change description when color is selected', () => {
    render(<ColorInputAtom name="color" description="Choisir une couleur" />);

    const input = screen.getByLabelText('Choisir une couleur');

    fireEvent.change(input, {
      target: { value: '#0000ff' },
    });

    expect(screen.getByText('#0000ff')).toBeInTheDocument();
  });
});
