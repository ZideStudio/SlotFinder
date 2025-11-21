import { getClassName } from '@Front/utils/getClassName';
import type { ComponentPropsWithRef } from 'react';

import './TextInputAtom.scss';

type TextInputAtomProps = ComponentPropsWithRef<'input'> & {
  classModifier?: string;
};

export const TextInputAtom = ({ className, classModifier, ...props }: TextInputAtomProps) => {
  const parentClassName = getClassName({
    defaultClassName: 'ds-text-input-atom',
    modifiers: classModifier,
    className,
  });

  return <input className={parentClassName} type="text" {...props} />;
};
