import { Heading } from "@Front/ui/atoms/Heading/Heading";
import { Button } from "@Front/ui/molecules/Button/Button";
import { useTranslation } from "react-i18next";
import logo from "../../../assets/svg/logo/colored_logo_no_bg.svg";
import { Icon } from "@Front/ui/atoms/Icon/Icon";

import "./Welcome.scss";

export const Welcome = () => {
  const { t } = useTranslation("welcome");

  return (
    <div className="welcome subgrid">
      <Button className="welcome__connexion-button">{t("connexion")}</Button>
      <Icon className="welcome__logo" icon={logo} />
      <Heading level={1}>{t("slotFinder")}</Heading>
      <Heading level={2}>{t("title")}</Heading>
      <Button className="welcome__event-button">{t("createEvent")}</Button>
    </div>
  );
};
