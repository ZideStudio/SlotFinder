import { Icon, type SvgIcon } from "@Front/ui/atoms/Icon/Icon";
import { getClassName } from "@Front/utils/helpers/getClassName";
import type { ComponentPropsWithoutRef, ElementType } from "react";

import "./ClickIcon.scss";

type ClickIconProps<Type extends ElementType = "button"> = {
  as?: Type;
  variant?: "default" | "bordered";
  icon: SvgIcon;
} & ComponentPropsWithoutRef<Type>;

export const ClickIcon = <Type extends ElementType = "button">({
  as,
  variant = "default",
  icon,
  className,
  ...props
}: ClickIconProps<Type>) => {
  const Component = as ?? "button";
  const isNativeButton = !as || as === "button";

  const parentClassName = getClassName({
    defaultClassName: "ds-click-icon",
    modifiers: [variant !== "default" && variant],
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
