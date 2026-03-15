import { render, screen } from '@testing-library/react';
import { NumberInput } from '../NumberInput';

describe('NumberInput', () => {
  it('should render a number input with label', () => {
    render(<NumberInput label="Test Label" name="test-input" />);

    expect(screen.getByRole('spinbutton', { name: 'Test Label' })).toBeInTheDocument();
  });

  it('should render required asterisk when required', () => {
    render(<NumberInput label="Test Label" name="test-input" required />);

    const asterisk = screen.getByText('*');
    expect(asterisk).toBeInTheDocument();
    expect(asterisk).toHaveAttribute('aria-hidden', 'true');
  });

  it('should apply custom id', () => {
    render(<NumberInput label="Test Label" name="test-input" id="custom-id" />);
    const input = screen.getByRole('spinbutton', { name: 'Test Label' });
    expect(input).toHaveAttribute('id', 'custom-id');
  });

  it('should apply custom className', () => {
    render(<NumberInput label="Test Label" name="test-input" className="custom-class" />);
    const input = screen.getByRole('spinbutton');

    const inputContainer = input.closest('div');
    expect(inputContainer).toHaveClass('ds-field custom-class');
  });

  it('should render error message linked with input', () => {
    render(<NumberInput label="Test Label" name="test-input" error="This is an error message" id="custom-test-id" />);

    const errorMessage = screen.getByText('This is an error message');
    expect(errorMessage).toBeInTheDocument();
    expect(errorMessage).toHaveAttribute('id', 'custom-test-id-error');

    const input = screen.getByRole('spinbutton');
    expect(input).toHaveAttribute('aria-invalid', 'true');
    expect(input).toHaveAttribute('aria-describedby', 'custom-test-id-error');
  });

  it('should remove aria-describedby when no error', () => {
    render(<NumberInput label="Test Label" name="test-input" id="custom-test-id" />);

    const input = screen.getByRole('spinbutton');
    expect(input).not.toHaveAttribute('aria-describedby');
  });
});
