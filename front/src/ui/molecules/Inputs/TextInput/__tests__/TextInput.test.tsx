import { render, screen } from '@testing-library/react';
import { TextInput } from '../TextInput';

describe('TextInputAtom', () => {
  it('should render a text input with label', () => {
    render(<TextInput label="Test Label" name="test-input" />);

    expect(screen.getByRole('textbox', { name: 'Test Label' })).toBeInTheDocument();
  });

  it('should render required asterisk when required', () => {
    render(<TextInput label="Test Label" name="test-input" required />);

    const asterisk = screen.getByText('*');
    expect(asterisk).toBeInTheDocument();
    expect(asterisk).toHaveAttribute('aria-hidden', 'true');
  });

  it('should apply custom id', () => {
    render(<TextInput label="Test Label" name="test-input" id="custom-id" />);
    const input = screen.getByRole('textbox', { name: 'Test Label' });
    expect(input).toHaveAttribute('id', 'custom-id');
  });

  it('should apply custom className', () => {
    render(<TextInput label="Test Label" name="test-input" className="custom-class" />);
    const input = screen.getByRole('textbox');

    const inputContainer = input.closest('div');
    expect(inputContainer).toHaveClass('ds-text-input custom-class');
  });

  it('should render error message linked with input', () => {
    render(<TextInput label="Test Label" name="test-input" error="This is an error message" id="custom-test-id" />);

    const errorMessage = screen.getByText('This is an error message');
    expect(errorMessage).toBeInTheDocument();
    expect(errorMessage).toHaveAttribute('id', 'custom-test-id-error');

    const input = screen.getByRole('textbox');
    expect(input).toHaveAttribute('aria-invalid', 'true');
    expect(input).toHaveAttribute('aria-describedby', 'custom-test-id-error');
  });

  it('should render remove aria-describedby when no error', () => {
    render(<TextInput label="Test Label" name="test-input" id="custom-test-id" />);

    const input = screen.getByRole('textbox');
    expect(input).not.toHaveAttribute('aria-describedby');
  });
});
