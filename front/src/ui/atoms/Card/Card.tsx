import { getClassName } from "@Front/utils/helpers/getClassName";
import type { ComponentPropsWithoutRef, ElementType } from "react";
import "./Card.scss";

type CardProps<As extends ElementType = "div"> = {
  as?: As;
  borderWeight?: "default" | "bold";
} & ComponentPropsWithoutRef<As>;

export const Card = <As extends ElementType = "div">({
  as,
  className,
  borderWeight = "default",
  children,
  ...props
}: CardProps<As>) => {
  const Component = as || "div";

  const parentClassName = getClassName({
    defaultClassName: "ds-card",
    className,
    modifiers: [borderWeight !== "default" && borderWeight],
  });

  return (
    <Component className={parentClassName} {...props}>
      {children}
    </Component>
  );
};
