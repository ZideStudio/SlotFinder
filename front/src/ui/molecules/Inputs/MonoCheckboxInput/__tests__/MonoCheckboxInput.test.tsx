import { render, screen } from '@testing-library/react';
import { MonoCheckboxInput } from '../MonoCheckboxInput';

describe('MonoCheckboxInput', () => {
  it('should render a checkbox input with label', () => {
    render(<MonoCheckboxInput label="Test Label" name="test-input" />);

    expect(screen.getByRole('checkbox', { name: 'Test Label' })).toBeInTheDocument();
  });

  it('should render required asterisk when required', () => {
    render(<MonoCheckboxInput label="Test Label" name="test-input" required />);

    const asterisk = screen.getByText('*');
    expect(asterisk).toBeInTheDocument();
    expect(asterisk).toHaveAttribute('aria-hidden', 'true');
  });

  it('should apply custom id', () => {
    render(<MonoCheckboxInput label="Test Label" name="test-input" id="custom-id" />);
    const input = screen.getByRole('checkbox', { name: 'Test Label' });
    expect(input).toHaveAttribute('id', 'custom-id');
  });

  it('should apply custom className', () => {
    render(<MonoCheckboxInput label="Test Label" name="test-input" className="custom-class" />);
    const input = screen.getByRole('checkbox');

    const inputContainer = input.closest('div');
    expect(inputContainer).toHaveClass('ds-field custom-class');
  });

  it('should render error message linked with input', () => {
    render(
      <MonoCheckboxInput label="Test Label" name="test-input" error="This is an error message" id="custom-test-id" />,
    );

    const errorMessage = screen.getByText('This is an error message');
    expect(errorMessage).toBeInTheDocument();
    expect(errorMessage).toHaveAttribute('id', 'custom-test-id-error');

    const input = screen.getByRole('checkbox');
    expect(input).toHaveAttribute('aria-invalid', 'true');
    expect(input).toHaveAttribute('aria-describedby', 'custom-test-id-error');
  });

  it('should render remove aria-describedby when no error', () => {
    render(<MonoCheckboxInput label="Test Label" name="test-input" id="custom-test-id" />);

    const input = screen.getByRole('checkbox');
    expect(input).not.toHaveAttribute('aria-describedby');
  });

  it('should have the is-checkbox modifier class when isCheckbox prop is true', () => {
    render(<MonoCheckboxInput label="Test Label" name="test-input" isCheckbox />);
    const input = screen.getByRole('checkbox');

    const inputContainer = input.closest('div');
    expect(inputContainer).toHaveClass('ds-field--is-checkbox');
  });
});
