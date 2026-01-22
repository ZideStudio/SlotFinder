import { getClassName } from '@Front/utils/getClassName';
import type { ComponentPropsWithRef, OptionHTMLAttributes } from 'react';

import './SelectInputAtom.scss';

type Option = {
  label: string;
  value: string | number;
  disabled?: boolean;
} & OptionHTMLAttributes<HTMLOptionElement>;

type SelectInputAtomProps = Omit<ComponentPropsWithRef<'select'>, 'name' | 'className'> & {
  name: string;
  options: Option[];
  className?: string;
  placeholder?: string;
};

export const SelectInputAtom = ({
  name,
  options,
  className,
  placeholder,
  ...props
}: SelectInputAtomProps) => {
  const parentClassName = getClassName({
    defaultClassName: 'ds-select-input-atom',
    className,
  });

  return (
      <select name={name} className={parentClassName} defaultValue={placeholder ? '' : undefined} {...props}>
        {placeholder && (
          <option value="" disabled>
            {placeholder}
          </option>
        )}
        {options.map((option) => (
          <option key={String(option?.value)} label={option.label} value={option.value} disabled={option.disabled} />
        ))}
      </select>
  );
};
