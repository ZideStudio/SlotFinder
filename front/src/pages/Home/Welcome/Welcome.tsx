import { appRoutes } from "@Front/routing/appRoutes";
import { Heading } from "@Front/ui/atoms/Heading/Heading";
import { useTranslation } from "react-i18next";
import { NavLink } from "react-router";

export const Welcome = () => {
  const { t } = useTranslation("welcome");

  return (
    <>
      <Heading level={1}>{t("title")}</Heading>
      <NavLink to={appRoutes.signUp()}>Sign Up</NavLink>
    </>
  );
};
