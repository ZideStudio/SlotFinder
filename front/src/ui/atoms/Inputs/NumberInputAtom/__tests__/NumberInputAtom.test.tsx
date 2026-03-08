import { render, screen } from '@testing-library/react';
import { NumberInputAtom } from '../NumberInputAtom';

describe('NumberInputAtom', () => {
  it('should render an input of type number with required name prop', () => {
    render(<NumberInputAtom name="test-input" placeholder="Enter text" />);
    const input = screen.getByRole('spinbutton');
    expect(input).toBeInTheDocument();
    expect(input).toHaveAttribute('type', 'number');
    expect(input).toHaveAttribute('name', 'test-input');
    expect(input).toHaveClass('ds-number-input-atom');
    expect(input).toHaveAttribute('placeholder', 'Enter text');
  });

  it('should apply custom className', () => {
    render(<NumberInputAtom name="test-input" className="custom-class" />);
    const input = screen.getByRole('spinbutton');
    expect(input).toHaveClass('ds-number-input-atom custom-class');
  });
});
