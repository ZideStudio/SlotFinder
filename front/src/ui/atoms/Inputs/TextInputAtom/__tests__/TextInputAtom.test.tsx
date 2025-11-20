import { render, screen } from '@testing-library/react';
import { TextInputAtom } from '../TextInputAtom';

describe('TextInputAtom', () => {
  it('renders an input of type text', () => {
    render(<TextInputAtom placeholder="Enter text" />);
    const input = screen.getByRole('textbox');
    expect(input).toBeInTheDocument();
    expect(input).toHaveAttribute('type', 'text');
    expect(input).toHaveClass('ds-text-input-atom');
    expect(input).toHaveAttribute('placeholder', 'Enter text');
  });

  it('applies classModifier as a modifier', () => {
    render(<TextInputAtom classModifier="large" />);
    const input = screen.getByRole('textbox');
    expect(input).toHaveClass('ds-text-input-atom ds-text-input-atom--large');
  });

  it('applies custom className', () => {
    render(<TextInputAtom className="custom-class" />);
    const input = screen.getByRole('textbox');
    expect(input).toHaveClass('ds-text-input-atom custom-class');
  });
});
