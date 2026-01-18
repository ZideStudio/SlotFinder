import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import { TextareaInput } from '../TextareaInput';

describe('TextareaInput', () => {
  it('associates the label with the textarea for accessibility', () => {
    const label = 'Description';
    render(<TextareaInput label={label} name="test-input" />);
    const textarea = screen.getByRole('textbox') as HTMLTextAreaElement;
    expect(textarea).toBeInTheDocument();
    expect(textarea.tagName.toLowerCase()).toBe('textarea');
  });

  it('uses the provided id for the textarea and label', () => {
    const id = 'custom-id';
    const label = 'Notes';
    render(<TextareaInput id={id} label={label} name="test-input" />);
    const textarea = screen.getByRole('textbox') as HTMLTextAreaElement;
    expect(textarea.id).toBe(id);
  });

  it('renders an error message and sets aria-describedby and aria-invalid when error is provided', () => {
    const id = 'err-id';
    const errorText = 'This field is required';
    const label = 'Comments';
    render(<TextareaInput id={id} label={label} error={errorText} name="test-input" />);
    const textarea = screen.getByRole('textbox') as HTMLTextAreaElement;
    const errorEl = screen.getByText(errorText);
    expect(errorEl).toBeInTheDocument();
    expect(errorEl.id).toBe(`${id}-error`);
    expect(textarea).toHaveAttribute('aria-describedby', `${id}-error`);
    expect(textarea).toHaveAttribute('aria-invalid', 'true');
  });

  it('applies the default and custom className to the container', () => {
    const custom = 'my-custom';
    const { container } = render(<TextareaInput label="L" className={custom} name="test-input" />);
    const wrapper = container.firstChild as HTMLElement | null;
    expect(wrapper).toBeTruthy();
    expect(wrapper).toHaveClass('ds-textarea-input');
    expect(wrapper).toHaveClass(custom);
  });
});
