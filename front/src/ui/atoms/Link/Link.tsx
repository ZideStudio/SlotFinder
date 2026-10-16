import { getClassName } from "@Front/utils/helpers/getClassName";
import type { ComponentPropsWithoutRef, ElementType } from "react";
import "./Link.scss";

type LinkProps<Type extends ElementType = "a"> = {
  as?: Type;
  openInNewTab?: boolean;
} & ComponentPropsWithoutRef<Type>;

export const Link = ({
  as,
  className,
  openInNewTab = false,
  children,
  ...props
}: LinkProps) => {
  const parentClassName = getClassName({
    defaultClassName: "ds-link",
    className,
  });
  const Component = as ?? "a";

  return (
    <Component
      className={parentClassName}
      {...(Component === "a" ? { target: openInNewTab ? "_blank" : undefined, rel: openInNewTab ? "noopener noreferrer" : undefined } : {})}
      {...props}
    >
      {children}
    </Component>
  );
};
