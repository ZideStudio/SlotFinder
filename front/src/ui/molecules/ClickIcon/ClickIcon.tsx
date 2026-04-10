// oxlint-disable react/jsx-props-no-spreading
import { Icon, type SvgIcon } from '@Front/ui/atoms/Icon/Icon';
import { getClassName } from '@Front/utils/getClassName';
import type { ComponentPropsWithoutRef } from 'react';

import './ClickIcon.scss';

type IconProps = {
  icon: SvgIcon;
} & ComponentPropsWithoutRef<'button'>;

export const ClickIcon = ({ icon, ...props }: IconProps) => {
  const parentClassName = getClassName({
    defaultClassName: 'ds-click-icon',
    className: props.className,
  });

  return (
    <button className={parentClassName} {...props}>
      <Icon icon={icon} />
    </button>
  );
};
