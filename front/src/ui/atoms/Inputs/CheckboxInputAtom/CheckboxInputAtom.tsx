import { getClassName } from '@Front/utils/getClassName';
import type { ComponentPropsWithRef } from 'react';

import './CheckboxInputAtom.scss';

export type CheckboxInputAtomProps = ComponentPropsWithRef<'input'>;

export const CheckboxInputAtom = ({ className, ...props }: CheckboxInputAtomProps) => {
  const parentClassName = getClassName({
    defaultClassName: 'ds-checkbox-input-atom',
    className,
  });

  return <input className={parentClassName} type="checkbox" {...props} />;
};
