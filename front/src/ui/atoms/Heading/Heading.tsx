// oxlint-disable react/jsx-props-no-spreading
import { getClassName } from '@Front/utils/getClassName';
import type { ReactNode } from 'react';
import './Heading.scss';

type HeadingTag = 'h1' | 'h2' | 'h3';
type HeadingProps = React.HTMLAttributes<HTMLHeadingElement> & {
  // oxlint-disable-next-line no-magic-numbers
  level: 1 | 2 | 3;
  children: ReactNode;
};

export const Heading = ({ level, className, ...props }: HeadingProps) => {
  const parentClassName = getClassName({
    defaultClassName: 'ds-heading',
    className,
    modifiers: [`level-${level}`],
  });

  const Tag = `h${level}` as HeadingTag;

  return <Tag {...props} className={parentClassName} />;
};
