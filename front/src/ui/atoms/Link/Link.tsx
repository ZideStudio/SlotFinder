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
  const openInNewTabProps = openInNewTab
    ? { target: "_blank", rel: "noopener noreferrer" }
    : {};

  return (
    <Component className={parentClassName} {...openInNewTabProps} {...props}>
      {children}
    </Component>
  );
};
