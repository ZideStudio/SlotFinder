import { render, screen } from '@testing-library/react';
import { DateInput } from '../DateInput';

describe('DateInputAtom', () => {
  it('should render a date input with label', () => {
    render(<DateInput label="Test Label" name="test-input" value="2026-01-01" />);

    const dateInput = screen.getByLabelText('Test Label');
    expect(dateInput).toBeInTheDocument();
    expect(dateInput).toHaveAttribute('type', 'date');
    expect(dateInput).toHaveAttribute('name', 'test-input');
    expect(dateInput).toHaveValue('2026-01-01');
  });

  it('should render required asterisk when required', () => {
    render(<DateInput label="Test Label" name="test-input" required />);

    const asterisk = screen.getByText('*');
    expect(asterisk).toBeInTheDocument();
    expect(asterisk).toHaveAttribute('aria-hidden', 'true');
  });

  it('should apply custom id', () => {
    render(<DateInput label="Test Label" name="test-input" id="custom-id" value="2026-01-01" />);
    const input = screen.getByDisplayValue('2026-01-01');
    expect(input).toHaveAttribute('id', 'custom-id');
  });

  it('should apply custom className', () => {
    render(<DateInput label="Test Label" name="test-input" className="custom-class" value="2026-01-01" />);
    const input = screen.getByDisplayValue('2026-01-01');

    const inputContainer = input.closest('div');
    expect(inputContainer).toHaveClass('ds-field custom-class');
  });

  it('should render error message linked with input', () => {
    render(
      <DateInput
        label="Test Label"
        name="test-input"
        error="This is an error message"
        id="custom-test-id"
        value="2026-01-01"
      />,
    );

    const errorMessage = screen.getByText('This is an error message');
    expect(errorMessage).toBeInTheDocument();
    expect(errorMessage).toHaveAttribute('id', 'custom-test-id-error');

    const input = screen.getByDisplayValue('2026-01-01');
    expect(input).toHaveAttribute('aria-invalid', 'true');
    expect(input).toHaveAttribute('aria-describedby', 'custom-test-id-error');
  });

  it('should render remove aria-describedby when no error', () => {
    render(<DateInput label="Test Label" name="test-input" id="custom-test-id" value="2026-01-01" />);

    const input = screen.getByDisplayValue('2026-01-01');
    expect(input).not.toHaveAttribute('aria-describedby');
  });
});
