import { getClassName } from '@Front/utils/getClassName';
import type { ComponentPropsWithRef } from 'react';

import './DateInputAtom.scss';

type DateInputAtomProps = Omit<ComponentPropsWithRef<'input'>, 'name'> & {
  name: string;
};

export const DateInputAtom = ({ className, ...props }: DateInputAtomProps) => {
  const parentClassName = getClassName({
    defaultClassName: 'ds-date-input-atom',
    className,
  });

  return <input className={parentClassName} type="date" {...props} />;
};
