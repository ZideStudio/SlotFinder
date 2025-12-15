import { render, screen } from '@testing-library/react';
import { TextInputAtom } from '../TextInputAtom';

describe('TextInputAtom', () => {
  it('should render an input of type text with required name prop', () => {
    render(<TextInputAtom name="test-input" placeholder="Enter text" />);
    const input = screen.getByRole('textbox');
    expect(input).toBeInTheDocument();
    expect(input).toHaveAttribute('type', 'text');
    expect(input).toHaveAttribute('name', 'test-input');
    expect(input).toHaveClass('ds-text-input-atom');
    expect(input).toHaveAttribute('placeholder', 'Enter text');
  });

  it('should apply custom className', () => {
    render(<TextInputAtom name="test-input" className="custom-class" />);
    const input = screen.getByRole('textbox');
    expect(input).toHaveClass('ds-text-input-atom custom-class');
  });
});
