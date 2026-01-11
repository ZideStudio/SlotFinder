import { getClassName } from '@Front/utils/getClassName';
import type { ComponentPropsWithRef } from 'react';

import './TextareaInputAtom.scss'

type TextareaInputAtomProps = Omit<ComponentPropsWithRef<'textarea'>, 'name'> & {
  name: string;
};

export const TextareaInputAtom = ({ className, ...props }: TextareaInputAtomProps) => {
  const parentClassName = getClassName({
    defaultClassName: 'ds-textarea-input-atom',
    className,
  });

  return <textarea className={parentClassName} {...props} />;
};
