import { getClassName } from '@Front/utils/getClassName';
import type { ComponentPropsWithRef, ChangeEvent } from 'react';

import './CheckboxInputAtom.scss';

export type CheckboxInputAtomProps = ComponentPropsWithRef<'input'> & {
  id: string;
  name?: string;
  onChange?: (e: ChangeEvent<HTMLInputElement>) => void;
  disabled?: boolean;
  required?: boolean;
};

export const CheckboxInputAtom = ({
  id,
  name,
  onChange,
  disabled,
  required,
  className,
  ...props
}: CheckboxInputAtomProps) => {
  const parentClassName = getClassName({
    defaultClassName: 'ds-checkbox-input-atom',
    className,
  });

  return (
    <input
      id={id}
      name={name}
      className={parentClassName}
      type="checkbox"
      onChange={onChange}
      disabled={disabled}
      required={required}
      {...props}
    />
  );
};
