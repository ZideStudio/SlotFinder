import { render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { SelectInputAtom } from '../SelectInputAtom';

describe('SelectInputAtom', () => {
  const options = [
    { label: 'One', value: '1' },
    { label: 'Two', value: '2', disabled: true },
  ];

  it('renders select with placeholder and options', () => {
    render(<SelectInputAtom id="test" name="test" options={options} placeholder="Select..." />);

    const select = screen.getByRole('combobox');
    expect(select).toBeInTheDocument();
    expect(select).toHaveAttribute('id', 'test');
    expect(select).toHaveAttribute('name', 'test');

    const placeholder = screen.getByRole('option', { name: 'Select...' });
    expect(placeholder).toBeInTheDocument();
    expect(placeholder).toHaveAttribute('value', '');
    expect(placeholder).toBeDisabled();
    expect(placeholder).toHaveProperty('selected', true);

    const renderedOptions = screen.getAllByRole('option');
    expect(renderedOptions).toHaveLength(options.length + 1);
  });

  it('renders select with options', () => {
    render(<SelectInputAtom id="test2" name="test2" options={options} />);

    options.forEach(option => {
      const optionElement = screen.getByRole('option', { name: option.label });
      expect(optionElement).toBeInTheDocument();
      expect(optionElement).toHaveAttribute('value', option.value);
    });
  });

  it('renders select with disabled option', () => {
    render(<SelectInputAtom id="test3" name="test3" options={options} />);
    
    const disabledOption = screen.getByRole('option', { name: 'Two' });
    expect(disabledOption).toBeDisabled();
  });

  it('applies className via getClassName', () => {
    render(<SelectInputAtom id="test2" name="test2" options={options} className="custom" />);

    const wrapper = screen.getByRole('combobox');
    expect(wrapper).toHaveClass('ds-select-input-atom');
    expect(wrapper).toHaveClass('custom');
    expect(wrapper).toBeInTheDocument();
  });
});
