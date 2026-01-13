import { getClassName } from '@Front/utils/getClassName';
import type { ComponentPropsWithRef } from 'react';
import './Link.scss';

type LinkProps = ComponentPropsWithRef<'a'> & {
  openInNewTab?: boolean;
};

export const Link = ({
  className,
  openInNewTab = false,
  ...props
}: LinkProps) => {
  const parentClassName = getClassName({
    defaultClassName: 'ds-link',
    className,
  });

  return (
    <a
      className={parentClassName}
      target={openInNewTab ? '_blank' : undefined}
      rel={openInNewTab ? 'noopener noreferrer' : undefined}
      {...props}
    />
  );
};
    
