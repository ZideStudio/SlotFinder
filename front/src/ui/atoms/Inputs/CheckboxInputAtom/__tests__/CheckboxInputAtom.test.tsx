import { render, screen } from '@testing-library/react';
import { CheckboxInputAtom } from '../CheckboxInputAtom';

describe('CheckboxInputAtom', () => {
  it('renders a checkbox input', () => {
    render(<CheckboxInputAtom name="test" />);
    const input = screen.getByRole('checkbox');
    expect(input).toBeInTheDocument();
    expect(input).toHaveAttribute('name', 'test');
    expect(input).toHaveClass('ds-checkbox-input-atom');
  });

  it('accepts additional class names', () => {
    render(<CheckboxInputAtom name="test" className="custom-class" />);
    const input = screen.getByRole('checkbox');
    expect(input).toHaveClass('custom-class');
  });
});
