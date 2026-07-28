import { appRoutes } from "@Front/routing/appRoutes";
import type { RouteHandle } from "@Front/routing/routeHandle";
import { Button } from "@Front/ui/molecules/Button/Button";
import { getClassName } from "@Front/utils/getClassName";
import AddCalendarIcon from "@material-symbols/svg-400/outlined/calendar_add_on.svg?react";
import Person from "@material-symbols/svg-400/outlined/person.svg?react";
import { useMemo } from "react";
import { NavLink, type UIMatch, useMatches } from "react-router";
import logo from "../../../../public/assets/logo.png";

import "./Header.scss";

type HeaderProps = {
  ignoreRouteHideHeader?: boolean;
  className?: string;
};

export const Header = ({
  ignoreRouteHideHeader = false,
  className = "",
}: HeaderProps) => {
  const matches = useMatches() as UIMatch<unknown, RouteHandle>[];
  const parentClassName = getClassName({
    defaultClassName: "header",
    className,
  });

  const hideHeader = useMemo(
    () =>
      ignoreRouteHideHeader
        ? false
        : matches.some((match) => match.handle?.hideHeader === true),
    [ignoreRouteHideHeader, matches],
  );

  if (hideHeader) {
    return null;
  }

  return (
    <header className={parentClassName}>
      <NavLink to={appRoutes.home()}>
        <img src={logo} alt="Slot Finder logo" className="header__logo" />
      </NavLink>
      <div className="header__buttons">
        <Button
          icon={AddCalendarIcon}
          variant="secondary"
          aria-label="add event"
          className="header__button"
        />
        <Button
          icon={Person}
          variant="secondary"
          aria-label="user profile"
          className="header__button"
        />
      </div>
    </header>
  );
};
