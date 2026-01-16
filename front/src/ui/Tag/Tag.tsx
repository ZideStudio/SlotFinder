import { getClassName } from '@Front/utils/getClassName';
import './Tag.scss';
import { type ReactNode } from 'react';

export type TagVariant ='default' |'success' | 'warning' | 'error';
export type TagAppearance = 'filled' | 'outlined';

export type TagProps = {
  variant?: TagVariant;
  className?: string;
  children: ReactNode;
  appearance?: TagAppearance;
};

export const Tag = ({ className, children, variant = 'default', appearance = 'filled' }: TagProps) => {
  const parentClassName = getClassName({
    defaultClassName: 'ds-tag',
    className,
    modifiers: [variant, appearance].filter(Boolean) as string[],
  });

  return (
    <span className={parentClassName}>
      <span className="ds-tag__content" title={typeof children === 'string' ? children : undefined}>
        {children}
      </span>
    </span>
  );
};
