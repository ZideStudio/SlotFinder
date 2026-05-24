import { fireEvent, render, screen } from '@testing-library/react';
import { useState } from 'react';
import { ColorInputAtom } from '../ColorInputAtom';

const ControlledColorInputAtom = ({ description }: { description: string }) => {
  const [value, setValue] = useState('');
  return (
    <ColorInputAtom name="color" description={description} value={value} onChange={e => setValue(e.target.value)} />
  );
};

describe('ColorInput', () => {
  it('renders the color input with default value', () => {
    render(<ColorInputAtom name="color" description="Choisir une couleur" value="" onChange={vi.fn()} />);

    expect(screen.getByText('Choisir une couleur')).toBeInTheDocument();

    const input = screen.getByLabelText('Choisir une couleur');
    expect(input).toBeInTheDocument();
    expect(input).toHaveValue('#000000');
  });

  it('updates the value when a new color is selected', () => {
    render(<ControlledColorInputAtom description="Choisir une couleur" />);

    const input = screen.getByLabelText('Choisir une couleur');

    fireEvent.change(input, {
      target: { value: '#00ff00' },
    });

    expect(input).toHaveValue('#00ff00');
  });

  it('applies extra props to the input', () => {
    render(<ColorInputAtom name="color" disabled description="La couleur" readOnly value="" onChange={vi.fn()} />);

    const input = screen.getByLabelText('La couleur');

    expect(input).toBeDisabled();
  });

  it('change description when color is selected', () => {
    render(<ControlledColorInputAtom description="Choisir une couleur" />);

    const input = screen.getByLabelText('Choisir une couleur');

    fireEvent.change(input, {
      target: { value: '#0000ff' },
    });

    expect(screen.getByText('#0000ff')).toBeInTheDocument();
  });
});
