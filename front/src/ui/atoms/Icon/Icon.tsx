import { ElementType, type FC, type SVGProps } from 'react';
import { getClassName } from '@Front/utils/getClassName';

import './Icon.scss';

export type SvgIcon = FC<SVGProps<SVGSVGElement>>;

export type IconProps = {
  icon: SvgIcon;
  className?: string;
};

export const Icon = ({ icon: IconComponent, className }: IconProps) => {
  const parentClassName = getClassName({
    defaultClassName: 'ds-icon',
    className,
  });

  return (
    <span className={parentClassName} aria-hidden="true" role="presentation">
      <IconComponent />
    </span>
  );
};
