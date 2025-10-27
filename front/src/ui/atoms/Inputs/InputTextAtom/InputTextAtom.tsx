import { getClassName } from '@Front/utils/getClassName';
import type { ComponentPropsWithRef } from 'react';

import './InputTextAtom.scss';

type InputTextAtomProps = ComponentPropsWithRef<'input'> & {
  classModifier?: string;
};

export const InputTextAtom = ({ className, classModifier, ...props }: InputTextAtomProps) => {
  const parentClassName = getClassName({
    defaultClassName: 'ds-input-text-atom',
    modifiers: classModifier,
    className,
  });

  return <input className={parentClassName} type="text" {...props} />;
};
