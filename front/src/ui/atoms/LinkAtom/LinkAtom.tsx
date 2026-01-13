import { getClassName } from '@Front/utils/getClassName';
import React, { ComponentPropsWithRef } from 'react';

type LinkAtomProps = ComponentPropsWithRef<'a'> & {
  openInNewTab?: boolean;
};

export const LinkAtom = ({
  className,
  openInNewTab = false,
  ...props
}: LinkAtomProps) => {
  const parentClassName = getClassName({
    defaultClassName: 'ds-link-atom',
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
    
