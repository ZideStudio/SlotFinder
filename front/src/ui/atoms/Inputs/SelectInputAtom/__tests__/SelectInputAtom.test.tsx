import { render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { SelectInputAtom } from '../SelectInputAtom';

describe('SelectInputAtom', () => {
  const options = [
    { label: 'One', value: '1' },
    { label: 'Two', value: '2', disabled: true },
  ];

  it('renders select with placeholder and options', () => {
    render(
      <SelectInputAtom id="test" name="test" options={options} placeholder="Select..." />
    );

    const select = screen.getByRole('combobox');
    expect(select).toBeInTheDocument();
    expect(select).toHaveAttribute('id', 'test');
    expect(select).toHaveAttribute('name', 'test');

    const placeholder = screen.getByText('Select...');
    expect(placeholder).toBeInTheDocument();
    expect(placeholder).toHaveAttribute('value', '');
    expect(placeholder).toBeDisabled();
    expect(placeholder).toHaveProperty('selected', true);

    const renderedOptions = screen.getAllByRole('option');
    expect(renderedOptions).toHaveLength(options.length + 1);

    const disabledOpt = screen.getByText('Two');
    expect(disabledOpt).toBeInTheDocument();
    expect(disabledOpt).toBeDisabled();
  });

  it('applies className via getClassName', () => {
    render(
      <SelectInputAtom id="test2" name="test2" options={options} className="custom" />
    );

    const wrapper = document.querySelector('.ds-select-input-atom.custom');
    expect(wrapper).toBeInTheDocument();
  });
});