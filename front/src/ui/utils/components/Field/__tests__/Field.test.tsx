import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import { Field } from '../Field';

const MockInput = (props: any) => <input {...props} />;

describe('Field Component', () => {
  it('Must link the label to the input via a generated ID', () => {
    render(<Field input={MockInput} label="Nom d'utilisateur" />);

    const label = screen.getByText("Nom d'utilisateur");
    const input = screen.getByRole('textbox', { name: "Nom d'utilisateur" });

    expect(label).toHaveAttribute('for', input.id);
    expect(input.id).toBeDefined();
  });

  it('Must use the manually provided ID instead of the generated one', () => {
    const manualId = 'custom-id';
    render(<Field input={MockInput} label="Email" id={manualId} />);

    const input = screen.getByRole('textbox', { name: 'Email' });
    expect(input.id).toBe(manualId);
  });

  it('Must display an error and configure the ARIA attributes', () => {
    const errorMessage = 'This field is required';
    render(<Field input={MockInput} label="Test" error={errorMessage} />);

    const input = screen.getByRole('textbox', { name: 'Test' });
    const errorElement = screen.getByText(errorMessage);

    expect(input).toHaveAttribute('aria-invalid', 'true');
    expect(input).toHaveAttribute('aria-describedby', errorElement.id);
  });

  it('Must not have aria-describedby if no error is present', () => {
    render(<Field input={MockInput} label="Test" />);
    const input = screen.getByRole('textbox', { name: 'Test' });

    expect(input).not.toHaveAttribute('aria-describedby');
    expect(input).toHaveAttribute('aria-invalid', 'false');
  });

  it('Must pass additional props to the input', () => {
    render(<Field input={MockInput} placeholder="Enter your text" name="username" label="Username" />);

    const input = screen.getByRole('textbox', { name: 'Username' });
    expect(input).toHaveAttribute('placeholder', 'Enter your text');
    expect(input).toHaveAttribute('name', 'username');
  });

  it('Must apply the CSS classes correctly', () => {
    const { container } = render(<Field input={MockInput} className="custom-class" />);

    const wrapper = container.firstChild;
    expect(wrapper).toHaveClass('ds-field');
    expect(wrapper).toHaveClass('custom-class');
  });
});
