import { getClassName } from '@Front/utils/getClassName';
import type { ComponentPropsWithRef, ChangeEvent } from 'react';

import './CheckboxInputAtom.scss';

export type CheckboxInputAtomProps = Omit<ComponentPropsWithRef<'input'>, 'name'> & {
  name: string;
};

export const CheckboxInputAtom = ({
  name,
  className,
  ...props
}: CheckboxInputAtomProps) => {
  const parentClassName = getClassName({
    defaultClassName: 'ds-checkbox-input-atom',
    className,
  });

  return (
    <input
      name={name}
      className={parentClassName}
      type="checkbox"
      {...props}
    />
  );
};
