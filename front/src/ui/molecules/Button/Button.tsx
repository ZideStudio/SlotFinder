import { Icon } from '@Front/ui/atoms/Icon/Icon';
import { Spinner } from '@Front/ui/atoms/Spinner/Spinner';
import { getClassName } from '@Front/utils/getClassName';
import type { ComponentPropsWithoutRef, ElementType, FC, SVGProps } from 'react';
import './Button.scss';

export type SvgIcon = FC<SVGProps<SVGSVGElement>>;

type ButtonProps<Type extends ElementType = 'button'> = {
  as?: Type;
  variant?: 'primary' | 'secondary';
  color?: 'default' | 'neutral' | 'danger';
  icon?: SvgIcon;
  disabled?: boolean;
  isLoading?: boolean;
} & ComponentPropsWithoutRef<Type>;

export const Button = <Type extends ElementType = 'button'>({
  as,
  className,
  children,
  variant = 'primary',
  color = 'default',
  icon,
  disabled,
  isLoading,
  ...props
}: ButtonProps<Type>) => {
  const Component = as ?? 'button';
  const isNativeButton = !as || as === 'button';
   if (isLoading) {
    disabled = true;
  }

  const parentClassName = getClassName({
    defaultClassName: 'ds-button',
    modifiers: [variant, color, disabled && 'disabled'],
    className,
  });

  return (
    <Component className={parentClassName} {...(isNativeButton && { type: 'button' })} {...props}>
      {isLoading && <Spinner />}
      {icon && <Icon icon={icon} />}
      {children}
    </Component>
  );
};
