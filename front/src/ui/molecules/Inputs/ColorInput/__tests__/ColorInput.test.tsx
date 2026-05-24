import { render, screen } from '@testing-library/react';
import { ColorInput } from '../ColorInput';

describe('ColorInput', () => {
  it('should render a color input with label', () => {
    render(
      <ColorInput label="Test Label" name="test-input" description="Choose a color" value="" onChange={vi.fn()} />,
    );
    expect(screen.getByLabelText('Test Label')).toBeInTheDocument();
  });

  it('should render required asterisk when required', () => {
    render(
      <ColorInput
        label="Test Label"
        name="test-input"
        description="Choose a color"
        required
        value=""
        onChange={vi.fn()}
      />,
    );
    const asterisk = screen.getByText('*');
    expect(asterisk).toBeInTheDocument();
    expect(asterisk).toHaveAttribute('aria-hidden', 'true');
  });

  it('should apply custom id', () => {
    render(
      <ColorInput
        label="Test Label"
        name="test-input"
        description="Choose a color"
        id="custom-id"
        value=""
        onChange={vi.fn()}
      />,
    );
    const input = screen.getByLabelText('Test Label');
    expect(input).toHaveAttribute('id', 'custom-id');
  });

  it('should apply custom className', () => {
    render(
      <ColorInput
        label="Test Label"
        name="test-input"
        description="Choose a color"
        className="custom-class"
        value=""
        onChange={vi.fn()}
      />,
    );
    const input = screen.getByLabelText('Test Label');
    const inputContainer = input.closest('div');
    expect(inputContainer).toHaveClass('ds-input custom-class');
  });

  it('should render error message linked with input', () => {
    render(
      <ColorInput
        label="Test Label"
        name="test-input"
        description="Choose a color"
        error="This is an error message"
        id="custom-test-id"
        value=""
        onChange={vi.fn()}
      />,
    );
    const errorMessage = screen.getByText('This is an error message');
    expect(errorMessage).toBeInTheDocument();
    expect(errorMessage).toHaveAttribute('id', 'custom-test-id-error');
    const input = screen.getByLabelText('Test Label');
    expect(input).toHaveAttribute('aria-invalid', 'true');
    expect(input).toHaveAttribute('aria-describedby', 'custom-test-id-error');
  });

  it('should remove aria-describedby when no error', () => {
    render(
      <ColorInput
        label="Test Label"
        name="test-input"
        description="Choose a color"
        id="custom-test-id"
        value=""
        onChange={vi.fn()}
      />,
    );
    const input = screen.getByLabelText('Test Label');
    expect(input).not.toHaveAttribute('aria-describedby');
  });
});
