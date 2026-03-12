import { getClassName } from '@Front/utils/getClassName';
import type { ComponentPropsWithRef } from 'react';

import './NumberInputAtom.scss';

type NumberInputAtomProps = Omit<ComponentPropsWithRef<'input'>, 'name'> & {
  name: string;
};

export const NumberInputAtom = ({ className, ...props }: NumberInputAtomProps) => {
  const parentClassName = getClassName({
    defaultClassName: 'ds-number-input-atom',
    className,
  });

  return <input className={parentClassName} type="number" {...props} />;
};
