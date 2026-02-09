import { getClassName } from '@Front/utils/getClassName';
import { getContrastTextColor } from '@Front/utils/getContrastTextColor';
import { type ReactNode } from 'react';

import './Tag.scss';

export type TagAppearance = 'filled' | 'outlined';

export type TagProps = {
  color: string;
  title?: string;
  className?: string;
  children: ReactNode;
  appearance?: TagAppearance;
};

export const Tag = ({ className, children, color, title, appearance = 'filled' }: TagProps) => {
  const parentClassName = getClassName({
    defaultClassName: 'ds-tag',
    className,
    modifiers: [appearance].filter(Boolean) as string[],
  });

  return (
    <span
      className={parentClassName}
      title={title}
      style={
        {
          '--tag-color': color,
          color: appearance === 'filled' ? getContrastTextColor(color) : 'black',
        } as React.CSSProperties
      }
    >
      {children}
    </span>
  );
};
