import { render, screen } from '@testing-library/react';
import { TextareaInputAtom } from '../TextareaInputAtom';

describe('TextareaInputAtom', () => {
  it('should render a textarea with required name prop', () => {
    render(<TextareaInputAtom name="test-textarea" placeholder="Enter text" />);
    const textarea = screen.getByRole('textbox');
    expect(textarea).toBeInTheDocument();
    expect(textarea).toHaveAttribute('name', 'test-textarea');
    expect(textarea).toHaveClass('ds-textarea-input-atom');
    expect(textarea).toHaveAttribute('placeholder', 'Enter text');
  });

  it('should apply custom className', () => {
    render(<TextareaInputAtom name="test-textarea" className="custom-class" />);
    const textarea = screen.getByRole('textbox');
    expect(textarea).toHaveClass('ds-textarea-input-atom custom-class');
  });
});
