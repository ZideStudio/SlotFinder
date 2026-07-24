import { Heading } from "@Front/ui/atoms/Heading/Heading";
import { useTranslation } from "react-i18next";
import { useLocation } from "react-router";

export const ErrorPage = () => {
  const { t } = useTranslation("error");
  const location = useLocation();

  return (
    <main>
      <Heading level={1}>{t("title")}</Heading>
      <p role="alert">{location.state?.message || t("unexpected")}</p>
    </main>
  );
};
