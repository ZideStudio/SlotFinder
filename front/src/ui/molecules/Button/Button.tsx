import type { ElementType, ComponentPropsWithoutRef, FC, SVGProps } from 'react';
import { getClassName } from '@Front/utils/getClassName';
import './Button.scss';
import { Icon } from '@Front/ui/atoms/Icon/Icon';

export type SvgIcon = FC<SVGProps<SVGSVGElement>>;

type ButtonProps<Type extends ElementType = 'button'> = {
  as?: Type;
  variant?: 'primary' | 'secondary';
  color?: 'default' | 'neutral' | 'danger';
  icon?: SvgIcon;
  disabled?: boolean;
} & ComponentPropsWithoutRef<Type>;

export const Button = <Type extends ElementType = 'button'>({ as, className, children, variant = 'primary', color = 'default', icon, disabled, ...props }: ButtonProps<Type>) => {
  const Component = as || 'button';

  const parentClassName = getClassName({
    defaultClassName: 'ds-button',
    modifiers: [
        variant,
        color,
        disabled && 'disabled',
    ],
    className,
  });

  return (
    <Component className={parentClassName} {...props}>
      {icon && <Icon icon={icon} />}
      {children}
    </Component>
  );
};
