import { getByLabelText, render } from '@testing-library/react';
import { FileUploadInputAtom } from '../FileUploadInputAtom';

describe('FileUploadInputAtom', () => {
  it('should render the component with custom class name', () => {
    const { container } = render(<FileUploadInputAtom name="file-upload" className="custom-class" />);
    const inputElement = getByLabelText(container, 'Déposer');
    const labelElement = inputElement.closest('label');
    expect(labelElement).toHaveClass('ds-file-upload-input-atom');
    expect(labelElement).toHaveClass('custom-class');
  });

  it('should render the component with generated id when no id is provided', () => {
    const { container } = render(<FileUploadInputAtom name="file-upload" />);
    const inputElement = getByLabelText(container, 'Déposer');
    expect(inputElement).toHaveAttribute('id');
  });

  it('should render the component with provided id', () => {
    const { container } = render(<FileUploadInputAtom name="file-upload" id="custom-id" />);
    const inputElement = getByLabelText(container, 'Déposer');
    expect(inputElement).toHaveAttribute('id', 'custom-id');
  });
});
