import { fireEvent, render, screen } from '@testing-library/react';
import { PictureUploadInput } from '../PictureUploadInput';

describe('PictureUploadInput', () => {
  beforeAll(() => {
    URL.createObjectURL = vi.fn(() => 'blob:http://localhost/fake-url');
    URL.revokeObjectURL = vi.fn();
  });

  afterAll(() => {
    vi.restoreAllMocks();
  });

  it('should render the component with label', () => {
    render(<PictureUploadInput label="Test Label" name="test-input" />);

    expect(screen.getByText('Test Label')).toBeInTheDocument();
  });

  it('should apply custom className', () => {
    render(<PictureUploadInput label="Test Label" name="test-input" className="custom-class" />);

    const container = screen.getByText('Test Label').closest('.ds-picture-upload-input');
    expect(container).toHaveClass('custom-class');
  });

  it('should render error message when error prop is provided', () => {
    render(<PictureUploadInput label="Test Label" name="test-input" error="This is an error message" />);

    expect(screen.getByText('This is an error message')).toBeInTheDocument();
  });

  it('should render image preview when a valid image file is selected', () => {
    render(<PictureUploadInput label="Test Label" name="test-input" />);

    const input = screen.getByLabelText('Test Label') as HTMLInputElement;
    const file = new File(['test'], 'test-image.png', { type: 'image/png' });

    fireEvent.change(input, { target: { files: [file] } });

    const img = screen.getByAltText('Preview') as HTMLImageElement;
    expect(img).toBeInTheDocument();
    expect(img).toHaveAttribute('src', 'blob:http://localhost/fake-url');
  });

  it('should not render image preview when a non-image file is selected', () => {
    render(<PictureUploadInput label="Test Label" name="test-input" />);

    const input = screen.getByLabelText('Test Label') as HTMLInputElement;
    const file = new File(['dummy content'], 'document.pdf', { type: 'application/pdf' });

    fireEvent.change(input, { target: { files: [file] } });

    expect(screen.queryByAltText('Preview')).not.toBeInTheDocument();
  });
});
