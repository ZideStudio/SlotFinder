import { render, screen } from '@testing-library/react';
import { SelectInput } from '../SelectInput';

describe('SelectInputAtom', () => {
  const options = [
    { label: 'One', value: '1' },
    { label: 'Two', value: '2', disabled: true },
  ];
  it('should render a Select input with label and options', () => {
    render(<SelectInput label="Test Label" name="test-input" options={options} />);

    expect(screen.getByRole('combobox', { name: 'Test Label' })).toBeInTheDocument();
  });

  it('should render a Select input with placeholder and options', () => {
    render(<SelectInput label="Test Label" name="test-input" options={options} placeholder="Select..." />);

    const placeholder = screen.getByRole('option', { name: 'Select...' });
    expect(placeholder).toBeInTheDocument();
    expect(placeholder).toHaveAttribute('value', '');
    expect(placeholder).toBeDisabled();
    expect(placeholder).toHaveProperty('selected', true);

    const renderedOptions = screen.getAllByRole('option');
    expect(renderedOptions).toHaveLength(options.length + 1);
  });

  it('should apply custom id', () => {
    render(<SelectInput label="Test Label" name="test-input" id="custom-id" options={options} />);
    const input = screen.getByRole('combobox', { name: 'Test Label' });
    expect(input).toHaveAttribute('id', 'custom-id');
  });

  it('should apply custom className', () => {
    render(<SelectInput label="Test Label" name="test-input" className="custom-class" options={options} />);
    const input = screen.getByRole('combobox');

    const inputContainer = input.closest('div');
    expect(inputContainer).toHaveClass('ds-field custom-class');
  });

  it('should render error message linked with input', () => {
    render(
      <SelectInput
        label="Test Label"
        name="test-input"
        error="This is an error message"
        id="custom-test-id"
        options={options}
      />,
    );

    const errorMessage = screen.getByText('This is an error message');
    expect(errorMessage).toBeInTheDocument();
    expect(errorMessage).toHaveAttribute('id', 'custom-test-id-error');

    const input = screen.getByRole('combobox');
    expect(input).toHaveAttribute('aria-invalid', 'true');
    expect(input).toHaveAttribute('aria-describedby', 'custom-test-id-error');
  });

  it('should render remove aria-describedby when no error', () => {
    render(<SelectInput label="Test Label" name="test-input" id="custom-test-id" options={options} />);

    const input = screen.getByRole('combobox');
    expect(input).not.toHaveAttribute('aria-describedby');
  });
});
