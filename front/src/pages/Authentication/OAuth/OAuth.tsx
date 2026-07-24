import { useTranslation } from "react-i18next";
import { useOAuth } from "./useOAuth";
import { Heading } from "@Front/ui/atoms/Heading/Heading";
import { Link } from "@Front/ui/atoms/Link/Link";

import "./OAuth.scss";

export const OAuth = () => {
  const { t } = useTranslation("authentication");
  const { oAuthProviders } = useOAuth();

  return (
    <nav aria-labelledby="oauth-provider-heading" className="oauth-nav">
      <Heading
        level={2}
        className="oauth-nav__heading"
        id="oauth-provider-heading"
      >
        {t("signInWithProvider")}
      </Heading>
      <ul className="oauth-nav__list">
        {oAuthProviders.map((provider) => (
          <li key={provider.label}>
            <Link
              href={provider.href}
              aria-label={t(provider.ariaLabel)}
              rel="noopener noreferrer"
            >
              {provider.icon}
              <span>{provider.label}</span>
            </Link>
          </li>
        ))}
      </ul>
    </nav>
  );
};
