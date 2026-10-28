import { Outlet, useMatches, type UIMatch } from "react-router";
import { useTranslation } from "react-i18next";
import type { AuthenticationRouteHandle } from "@Front/routing/routeHandle";

import { CardPage } from "@Front/components/CardPage/CardPage";
import { ExistingAccount } from "./ExistingAccount/ExistingAccount";
import { OAuth } from "./OAuth/OAuth";
import "./Authentication.scss";

export const Authentication = () => {
  const { t } = useTranslation();
  const matches = useMatches() as UIMatch<unknown, AuthenticationRouteHandle>[];
  const currentTitle = matches.at(-1)?.handle?.title;

  return (
    <CardPage title={currentTitle ? t(currentTitle) : ""}>
      <Outlet />
      <hr className="authentication-divider" />
      <OAuth />
      <ExistingAccount />
    </CardPage>
  );
};
