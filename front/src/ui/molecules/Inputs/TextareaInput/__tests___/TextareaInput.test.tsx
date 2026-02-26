import { render, screen } from '@testing-library/react';
import { TextareaInput } from '../TextareaInput';

describe('TextareaInput', () => {
  it('should render a textarea input with label', () => {
    render(<TextareaInput label="Test Label" name="test-input" />);

    expect(screen.getByRole('textbox', { name: 'Test Label' })).toBeInTheDocument();
  });

  it('should render required asterisk when required', () => {
    render(<TextareaInput label="Test Label" name="test-input" required />);

    const asterisk = screen.getByText('*');
    expect(asterisk).toBeInTheDocument();
    expect(asterisk).toHaveAttribute('aria-hidden', 'true');
  });

  it('should apply custom id', () => {
    render(<TextareaInput label="Test Label" name="test-input" id="custom-id" />);
    const input = screen.getByRole('textbox', { name: 'Test Label' });
    expect(input).toHaveAttribute('id', 'custom-id');
  });

  it('should apply custom className', () => {
    render(<TextareaInput label="Test Label" name="test-input" className="custom-class" />);
    const input = screen.getByRole('textbox');

    const inputContainer = input.closest('div');
    expect(inputContainer).toHaveClass('ds-field custom-class');
  });

  it('should render error message linked with input', () => {
    render(<TextareaInput label="Test Label" name="test-input" error="This is an error message" id="custom-test-id" />);

    const errorMessage = screen.getByText('This is an error message');
    expect(errorMessage).toBeInTheDocument();
    expect(errorMessage).toHaveAttribute('id', 'custom-test-id-error');

    const input = screen.getByRole('textbox');
    expect(input).toHaveAttribute('aria-invalid', 'true');
    expect(input).toHaveAttribute('aria-describedby', 'custom-test-id-error');
  });

  it('should render remove aria-describedby when no error', () => {
    render(<TextareaInput label="Test Label" name="test-input" id="custom-test-id" />);

    const input = screen.getByRole('textbox');
    expect(input).not.toHaveAttribute('aria-describedby');
  });
});
