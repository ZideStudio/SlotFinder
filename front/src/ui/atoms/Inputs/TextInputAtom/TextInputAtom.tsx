import { getClassName } from '@Front/utils/getClassName';
import type { ComponentPropsWithRef } from 'react';

import './TextInputAtom.scss';

type TextInputAtomProps = ComponentPropsWithRef<'input'> & {
  name: string;
};

export const TextInputAtom = ({ className, ...props }: TextInputAtomProps) => {
  const parentClassName = getClassName({
    defaultClassName: 'ds-text-input-atom',
    className,
  });

  return <input className={parentClassName} type="text" {...props} />;
};
