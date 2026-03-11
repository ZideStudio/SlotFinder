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
  const today = new Date().toISOString().split('T')[0];

  return <input className={parentClassName} type="date" defaultValue={today} {...props} />;
};
