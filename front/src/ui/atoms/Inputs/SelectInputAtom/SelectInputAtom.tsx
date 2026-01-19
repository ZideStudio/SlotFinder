import { getClassName } from '@Front/utils/getClassName';
import type { ComponentPropsWithRef } from 'react';

import './SelectInputAtom.scss';

type Option = {
  label: string;
  value: string | number;
  disabled?: boolean;
};

type SelectInputAtomProps = Omit<ComponentPropsWithRef<'select'>, 'name' | 'id' | 'className' | 'required'> & {
  id?: string;
  name: string;
  options: Option[];
  className?: string;
  placeholder?: string;
};

export const SelectInputAtom = ({
  id,
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
      <select id={id} name={name} className={parentClassName} defaultValue={placeholder ? '' : undefined} {...props}>
        {placeholder && (
          <option value="" disabled>
            {placeholder}
          </option>
        )}
        {options.map(opt => (
          <option key={String(opt.value)} value={opt.value} disabled={opt.disabled}>
            {opt.label}
          </option>
        ))}
      </select>
  );
};
