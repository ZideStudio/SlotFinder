import { Icon, type SvgIcon } from "@Front/ui/atoms/Icon/Icon";
import { getClassName } from "@Front/utils/getClassName";
import type { ComponentPropsWithoutRef, ElementType } from "react";

import "./ClickIcon.scss";

type ClickIconProps<Type extends ElementType = "button"> = {
  as?: Type;
  icon: SvgIcon;
} & ComponentPropsWithoutRef<Type>;

export const ClickIcon = <Type extends ElementType = "button">({
  as,
  icon,
  className,
  ...props
}: ClickIconProps<Type>) => {
  const Component = as ?? "button";
  const isNativeButton = !as || as === "button";

  const parentClassName = getClassName({
    defaultClassName: "ds-click-icon",
    className,
  });

  return (
    <Component
      className={parentClassName}
      {...(isNativeButton && { type: "button" })}
      {...props}
    >
      <Icon icon={icon} />
    </Component>
  );
};
