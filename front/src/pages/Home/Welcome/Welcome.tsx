import { Heading } from "@Front/ui/atoms/Heading/Heading";
import { Button } from "@Front/ui/molecules/Button/Button";
import { useTranslation } from "react-i18next";
import logo from "../../../assets/svg/logo/colored_logo_no_bg.svg";
import { Icon } from "@Front/ui/atoms/Icon/Icon";

import "./Welcome.scss";

export const Welcome = () => {
  const { t } = useTranslation("welcome");

  return (
    <section className="welcome subgrid">
      <Button className="welcome__sign-in-button">{t("signIn")}</Button>

      <Heading level={1} className="welcome__slot-finder">
        {t("slotFinder")}
      </Heading>
      <Heading level={2} className="welcome__tag-line">
        {t("tagLine")}
      </Heading>

      <Icon className="welcome__logo" icon={logo} />

      <Button className="welcome__event-button">{t("createEvent")}</Button>
    </section>
  );
};
