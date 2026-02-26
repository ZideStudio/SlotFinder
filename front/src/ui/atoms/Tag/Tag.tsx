import { getClassName } from '@Front/utils/getClassName';
import { getContrastTextColor } from '@Front/utils/getContrastTextColor';
import { type CSSProperties, type ReactNode } from 'react';

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

const tagTitle = title || (typeof children === 'string' ? children : undefined);

  return (
    <span
      className={parentClassName}
      title={tagTitle}
      style={
        {
          '--tag-color': color,
          color: appearance === 'filled' ? getContrastTextColor(color) : 'black',
        } as CSSProperties
      }
    >
      {children}
    </span>
  );
};
