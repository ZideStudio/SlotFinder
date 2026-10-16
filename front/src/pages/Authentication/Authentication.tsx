import { Outlet, useMatches, type UIMatch } from "react-router";
import { useTranslation } from "react-i18next";
import type { AuthenticationRouteHandle } from "@Front/routing/routeHandle";

import "./Authentication.scss";
import { CardPage } from "@Front/components/CardPage/CardPage";

export const Authentication = () => {
  const { t } = useTranslation();
  const matches = useMatches() as UIMatch<unknown, AuthenticationRouteHandle>[];
  const currentTitle = matches.at(-1)?.handle?.title;

  return (
    <section className="authentication subgrid">
      <CardPage
        title={currentTitle ? t(currentTitle) : ""}
        className="authentication__content"
      >
        <Outlet />
      </CardPage>
    </section>
  );
};
