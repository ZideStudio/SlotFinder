import { ElementType, ComponentPropsWithoutRef, ReactNode } from 'react';
import { getClassName } from '@Front/utils/getClassName';
import './Card.scss';

type CardProps<T extends ElementType = 'div'> = {
  as?: T;
} & ComponentPropsWithoutRef<T>;

export const Card = <T extends ElementType = 'div'>({ as, className, children, ...props }: CardProps<T>) => {
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
