import Logo from "@Front/assets/svg/logo/black_logo_no_bg.svg";
import { appRoutes } from "@Front/routing/appRoutes";
import { Card } from "@Front/ui/atoms/Card/Card";
import { Heading } from "@Front/ui/atoms/Heading/Heading";
import { Icon } from "@Front/ui/atoms/Icon/Icon";
import { getClassName } from "@Front/utils/getClassName";
import type { ComponentProps, ReactNode } from "react";
import { NavLink } from "react-router";

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
        <NavLink
          to={appRoutes.home()}
          className="card-page__logo"
          aria-label="Slot Finder home page"
        >
          <Icon icon={Logo} />
        </NavLink>
      )}
      <Heading level={1} className={`${defaultClassName}__title`}>
        {title}
      </Heading>
      <Card borderWeight="bold">{children}</Card>
    </section>
  );
};
