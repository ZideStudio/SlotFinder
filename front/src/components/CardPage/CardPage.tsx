import { appRoutes } from "@Front/routing/appRoutes";
import { Card } from "@Front/ui/atoms/Card/Card";
import { Heading } from "@Front/ui/atoms/Heading/Heading";
import { getClassName } from "@Front/utils/getClassName";
import type { ComponentProps, ReactNode } from "react";
import { NavLink } from "react-router";
import logo from "../../../public/assets/black_logo_without_background.png";

import "./CardPage.scss";

type CardProps = {
  title: ReactNode;
  isIconDisplayed?: boolean;
  children: ReactNode;
} & ComponentProps<"section">;

const defaultClassName = "card-page";

export const CardPage = ({
  title,
  isIconDisplayed = true,
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
      {isIconDisplayed && (
        <NavLink to={appRoutes.home()} className="card-page__logo">
          <img src={logo} alt="Slot Finder logo" />
        </NavLink>
      )}
      <Heading level={1} className={`${defaultClassName}__title`}>
        {title}
      </Heading>
      <Card borderWeight="bold">{children}</Card>
    </section>
  );
};
