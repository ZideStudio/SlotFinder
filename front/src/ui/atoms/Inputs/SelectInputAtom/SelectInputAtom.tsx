import { getClassName } from '@Front/utils/getClassName';
import type { ComponentPropsWithRef } from 'react';

import './SelectInputAtom.scss';

type Option = {
  label: string;
  value: string | number;
  disabled?: boolean;
};

type SelectInputAtomProps = Omit<ComponentPropsWithRef<'select'>, 'name'> & {
  id: string;
  name: string;
  options: Option[];
  error?: string;
  required?: boolean;
  className?: string;
  placeholder?: string;
};

export const SelectInputAtom = ({
  id,
  name,
  options,
  error,
  required,
  className,
  placeholder,
  ...props
}: SelectInputAtomProps) => {
  const parentClassName = getClassName({
    defaultClassName: 'ds-select-input-atom',
    className,
  });

  return (
    <div className={parentClassName}>
      <select id={id} name={name} className="ds-select-input-atom__select" {...props}>
        {placeholder && (
          <option value="" disabled selected>
            {placeholder}
          </option>
        )}
        {options.map(opt => (
          <option key={String(opt.value)} value={opt.value} disabled={opt.disabled}>
            {opt.label}
          </option>
        ))}
      </select>
    </div>
  );
};
