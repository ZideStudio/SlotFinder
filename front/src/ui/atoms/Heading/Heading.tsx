import { getClassName } from '@Front/utils/getClassName';
import type { HTMLAttributes, ReactNode } from 'react';
import './Heading.scss';

type HeadingTag = 'h1' | 'h2' | 'h3';
type HeadingProps = HTMLAttributes<HTMLHeadingElement> & {
  level: 1 | 2 | 3;
  children: ReactNode;
};

export const Heading = ({ level, className, ...props }: HeadingProps) => {
  const parentClassName = getClassName({
    defaultClassName: 'ds-heading',
    className,
    modifiers: [`level-${level}`],
  });

  const Tag = `h${level}` satisfies HeadingTag;

  return <Tag {...props} className={parentClassName} />;
};
