import { Card } from "@Front/ui/atoms/Card/Card";
import { Heading } from "@Front/ui/atoms/Heading/Heading";
import { getClassName } from "@Front/utils/getClassName";
import type { ComponentProps, ReactNode } from "react";

import "./CardPage.scss";

type CardProps = {
  title: ReactNode;
  children: ReactNode;
} & ComponentProps<"section">;

const defaultClassName = "card-page";

export const CardPage = ({
  title,
  className,
  children,
  ...props
}: CardProps) => {
  const parentClassName = getClassName({
    defaultClassName,
    className,
  });

  return (
    <section className={parentClassName} {...props}>
      <Heading level={1} className={`${defaultClassName}__title`}>
        {title}
      </Heading>
      <Card borderWeight="bold">{children}</Card>
    </section>
  );
};
