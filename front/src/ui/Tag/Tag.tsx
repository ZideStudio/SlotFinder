import { getClassName } from '@Front/utils/getClassName';
import './Tag.scss';
import { type ReactNode } from 'react';

export type TagAppearance = 'filled' | 'outlined';

export type TagProps = {
  color: string;
  className?: string;
  children: ReactNode;
  appearance?: TagAppearance;
};

export const Tag = ({ className, children, color, appearance = 'filled' }: TagProps) => {
  const parentClassName = getClassName({
    defaultClassName: 'ds-tag',
    className,
    modifiers: [appearance].filter(Boolean) as string[],
  });

  return (
    <span
      className={parentClassName}
      style={{ backgroundColor: appearance === 'filled' ? color : 'transparent', borderColor: color }}
    >
      <span className="ds-tag__content" title={typeof children === 'string' ? children : undefined}>
        {children}
      </span>
    </span>
  );
};
