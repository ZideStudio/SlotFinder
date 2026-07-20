import { useTranslation } from "react-i18next";

import { useOAuth } from "./useOAuth";
import { Heading } from "@Front/ui/atoms/Heading/Heading";
import { Link } from "@Front/ui/atoms/Link/Link";

import "./OAuth.css";

export const OAuth = () => {
  const { t } = useTranslation("authentication");
  const { oAuthProviders } = useOAuth();

  return (
    <Grid
      component="nav"
      container
      colSpan={{ "desktop-small": 4, tablet: 4, mobile: 4 }}
      colStart={{ "desktop-small": 5, tablet: 3, mobile: 1 }}
      aria-labelledby="oauth-provider-heading"
      className="oauth-nav"
    >
      <Heading level={2} id="oauth-provider-heading">
        {t("signInWithProvider")}
      </Heading>
      <ul className="oauth-nav__list">
        {oAuthProviders.map((provider) => (
          <li key={provider.label}>
            <Link
              href={provider.href}
              aria-label={t(provider.ariaLabel)}
              openInNewTab
            >
              {provider.icon}
              <span>{provider.label}</span>
            </Link>
          </li>
        ))}
      </ul>
    </Grid>
  );
};
