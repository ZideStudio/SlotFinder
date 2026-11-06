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
  const currentHandle = matches.at(-1)?.handle;
  const currentTitle = currentHandle?.title;
  const shouldDisplayOAuth = !currentHandle?.hideOAuth;

  return (
    <CardPage title={currentTitle ? t(currentTitle) : ""}>
      <Outlet />
      {shouldDisplayOAuth && (
        <>
          <hr className="authentication-divider" />
          <OAuth />
        </>
      )}
      <ExistingAccount />
    </CardPage>
  );
};
