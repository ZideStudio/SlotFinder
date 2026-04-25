import { render, screen } from '@testing-library/react';
import { PictureUploadInput } from '../PictureUploadInput';
import userEvent from '@testing-library/user-event';

describe('PictureUploadInput', () => {
  beforeAll(() => {
    URL.createObjectURL = vi.spyOn(URL, 'createObjectURL').mockImplementation(() => 'blob:http://localhost/fake-url');
    URL.revokeObjectURL = vi.spyOn(URL, 'revokeObjectURL').mockImplementation(() => undefined);
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

  it('should render image preview when a valid image file is selected', async () => {
  const user = userEvent.setup();
  render(<PictureUploadInput label="Test Label" name="test-input" />);

  const input = screen.getByLabelText('Test Label');
  const file = new File(['test'], 'test-image.png', { type: 'image/png' });

  await user.upload(input, file);

  const img = screen.getByAltText('Preview');
  expect(img).toBeInTheDocument();
  expect(img).toHaveAttribute('src', 'blob:http://localhost/fake-url');
});

it('should not render image preview when a non-image file is selected', async () => {
  const user = userEvent.setup();
  render(<PictureUploadInput label="Test Label" name="test-input" />);

  const input = screen.getByLabelText('Test Label');
  const file = new File(['dummy content'], 'document.pdf', { type: 'application/pdf' });

  await user.upload(input, file);

  expect(screen.queryByAltText('Preview')).toHaveAttribute('hidden');
});
});
