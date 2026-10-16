import { Outlet, useMatches, type UIMatch } from "react-router";
import { useTranslation } from "react-i18next";
import type { AuthenticationRouteHandle } from "@Front/routing/routeHandle";

import { CardPage } from "@Front/components/CardPage/CardPage";
import { OAuth } from "./OAuth/OAuth";

export const Authentication = () => {
  const { t } = useTranslation();
  const matches = useMatches() as UIMatch<unknown, AuthenticationRouteHandle>[];
  const currentTitle = matches.at(-1)?.handle?.title;

  return (
    <CardPage title={currentTitle ? t(currentTitle) : ""}>
      <Outlet />
      <hr />
      <OAuth />
    </CardPage>
  );
};
