import { render, screen } from "@testing-library/react";
import { DateInputAtom } from "../DateInputAtom";

describe('DateInputAtom', () => {
  it('renders the date input with default value', () => {
    render(<DateInputAtom name="date-input" value="2026-01-01" data-testid="date-input" />);
    const input = screen.getByTestId('date-input');
    expect(input).toBeInTheDocument();
    expect(input).toHaveAttribute('value', '2026-01-01');
    expect(input).toHaveAttribute('type', 'date');
    expect(input).toHaveAttribute('name', 'date-input');
    expect(input).toHaveClass('ds-date-input-atom');
  });

  it('should apply custom className', () => {
    render(<DateInputAtom name="date-input" className="custom-class" data-testid="date-input" />);
    const input = screen.getByTestId('date-input');
    expect(input).toHaveClass('ds-date-input-atom custom-class');
  });
});