import { render, screen } from '@testing-library/react';
import type { ComponentProps } from 'react';
import { Input } from '../Input';

const MockInput = (props: ComponentProps<'input'>) => <input {...props} />;

describe('Input', () => {
  it('should link the label to the input via a generated ID', () => {
    render(<Input input={MockInput} label="Nom d'utilisateur" />);

    const input = screen.getByRole('textbox', { name: "Nom d'utilisateur" });
    const label = screen.getByText("Nom d'utilisateur");

    expect(label).toHaveAttribute('for', input.id);
    expect(input.id).toBeDefined();
  });

  it('should use the manually provided ID instead of the generated one', () => {
    const manualId = 'custom-id';
    render(<Input input={MockInput} label="Email" id={manualId} />);

    const input = screen.getByRole('textbox', { name: 'Email' });
    expect(input.id).toBe(manualId);
  });

  it('should display an error and configure the ARIA attributes', () => {
    const errorMessage = 'This field is required';
    render(<Input input={MockInput} label="Test" error={errorMessage} />);

    const input = screen.getByRole('textbox', { name: 'Test' });
    const errorElement = screen.getByText(errorMessage);

    expect(input).toHaveAttribute('aria-invalid', 'true');
    expect(input).toHaveAttribute('aria-describedby', errorElement.id);
  });

  it('should not have aria-describedby if no error is present', () => {
    render(<Input input={MockInput} label="Test" />);
    const input = screen.getByRole('textbox', { name: 'Test' });

    expect(input).not.toHaveAttribute('aria-describedby');
    expect(input).toHaveAttribute('aria-invalid', 'false');
  });

  it('should pass additional props to the input', () => {
    render(<Input input={MockInput} placeholder="Enter your text" name="username" label="Username" />);

    const input = screen.getByRole('textbox', { name: 'Username' });
    expect(input).toHaveAttribute('placeholder', 'Enter your text');
    expect(input).toHaveAttribute('name', 'username');
  });

  it('should apply the CSS classes correctly', () => {
    const { container } = render(<Input input={MockInput} label="Test" className="custom-class" />);

    const wrapper = container.firstChild;
    expect(wrapper).toHaveClass('ds-input');
    expect(wrapper).toHaveClass('custom-class');
  });
});
