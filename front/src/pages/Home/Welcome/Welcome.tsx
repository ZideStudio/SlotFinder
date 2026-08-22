import { appRoutes } from "@Front/routing/appRoutes";
import { Heading } from "@Front/ui/atoms/Heading/Heading";
import { Button } from "@Front/ui/molecules/Button/Button";
import { useTranslation } from "react-i18next";
import { NavLink } from "react-router";

import './Welcome.scss';

export const Welcome = () => {
  const { t } = useTranslation("welcome");

  return (
    <div className="welcome">
      <Button className="welcome__connexion-button">{t("connexion")}</Button>
      <Heading level={1}>{t("slotFinder")}</Heading>

      <NavLink to={appRoutes.signUp()}>Sign Up</NavLink>
    </div>
  );
};
