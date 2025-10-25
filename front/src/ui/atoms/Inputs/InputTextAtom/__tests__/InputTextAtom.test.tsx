import { render, screen } from '@testing-library/react';
import { InputTextAtom } from '../InputTextAtom';

describe('InputTextAtom', () => {
  it('renders an input of type text', () => {
    render(<InputTextAtom placeholder="Enter text" />);
    const input = screen.getByRole('textbox');
    expect(input).toBeInTheDocument();
    expect(input).toHaveAttribute('type', 'text');
    expect(input).toHaveClass('ds-input-text-atom');
    expect(input).toHaveAttribute('placeholder', 'Enter text');
  });

  it('applies classModifier as a modifier', () => {
    render(<InputTextAtom classModifier="large" />);
    const input = screen.getByRole('textbox');
    expect(input).toHaveClass('ds-input-text-atom ds-input-text-atom--large');
  });

  it('applies custom className', () => {
    render(<InputTextAtom className="custom-class" />);
    const input = screen.getByRole('textbox');
    expect(input).toHaveClass('ds-input-text-atom custom-class');
  });
});
