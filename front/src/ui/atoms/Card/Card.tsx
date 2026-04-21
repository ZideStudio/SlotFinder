// oxlint-disable react/jsx-props-no-spreading
import { getClassName } from '@Front/utils/getClassName';
import type { ComponentPropsWithoutRef, ElementType } from 'react';
import './Card.scss';

type CardProps<As extends ElementType = 'div'> = {
  as?: As;
} & ComponentPropsWithoutRef<As>;

export const Card = <As extends ElementType = 'div'>({ as, className, children, ...props }: CardProps<As>) => {
  const Component = as || 'div';

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
