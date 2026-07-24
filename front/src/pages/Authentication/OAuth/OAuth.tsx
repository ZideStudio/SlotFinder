import { useTranslation } from "react-i18next";
import "./OAuth.scss";
import { useOAuth } from "./useOAuth";

export const OAuth = () => {
  const { t } = useTranslation("authentication");
  const { oAuthProviders } = useOAuth();

  return (
    <nav aria-labelledby="oauth-provider-heading" className="oauth-nav">
      <h2 className="oauth-nav__heading" id="oauth-provider-heading">
        {t("signInWithProvider")}
      </h2>
      <ul className="oauth-nav__list">
        {oAuthProviders.map((provider) => (
          <li key={provider.label}>
            <a
              href={provider.href}
              aria-label={t(provider.ariaLabel)}
              rel="noopener noreferrer"
            >
              {provider.icon}
              <span>{provider.label}</span>
            </a>
          </li>
        ))}
      </ul>
    </nav>
  );
};
