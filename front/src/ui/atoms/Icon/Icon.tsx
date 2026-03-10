import { type FC, type SVGProps } from 'react';
import { getClassName } from '@Front/utils/getClassName';

import './Icon.scss';

export type SvgIcon = FC<SVGProps<SVGSVGElement>>;

export type IconProps = SVGProps<SVGSVGElement> & {
  icon: SvgIcon;
};

export const Icon = ({ icon: IconComponent, className, ...props }: IconProps) => {
  const parentClassName = getClassName({
    defaultClassName: 'ds-icon',
    className,
  });

  return <IconComponent className={parentClassName} aria-hidden="true" role="presentation" {...props} />;
};
