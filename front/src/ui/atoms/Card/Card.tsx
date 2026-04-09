// oxlint-disable react/jsx-props-no-spreading
// oxlint-disable id-length
import { getClassName } from '@Front/utils/getClassName';
import type { ComponentPropsWithoutRef, ElementType } from 'react';
import './Card.scss';

type CardProps<T extends ElementType = 'div'> = {
  as?: T;
} & ComponentPropsWithoutRef<T>;

export const Card = <T extends ElementType = 'div'>({ as, className, children, ...props }: CardProps<T>) => {
  const Component = as ?? 'div';

  const parentClassName = getClassName({
    defaultClassName: 'ds-card',
    className,
  });

  return (
    <Component className={parentClassName} {...props}>
      {children}
    </Component>
  );
};
