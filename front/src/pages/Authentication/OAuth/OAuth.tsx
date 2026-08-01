import { useTranslation } from "react-i18next";
import { useOAuth } from "./useOAuth";
import { ClickIcon } from "@Front/ui/molecules/ClickIcon/ClickIcon";

import "./OAuth.scss";

export const OAuth = () => {
  const { t } = useTranslation("authentication");
  const { oAuthProviders } = useOAuth();

  return (
    <nav aria-label={t("signInWithProvider")} className="oauth-nav">
      <ul className="oauth-nav__list">
        {oAuthProviders.map((provider) => (
          <li key={provider.label}>
            <ClickIcon
              as="a"
              variant="bordered"
              href={provider.href}
              aria-label={t(provider.ariaLabel)}
              rel="noopener noreferrer"
              icon={provider.icon}
            />
          </li>
        ))}
      </ul>
    </nav>
  );
};
