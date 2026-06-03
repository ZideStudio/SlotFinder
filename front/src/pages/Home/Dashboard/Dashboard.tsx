import { Header } from "@Front/components/Layout/Header/Header";
import { useTranslation } from "react-i18next";

export const Dashboard = () => {
  <Header />;
  const { t } = useTranslation("dashboard");

  return <h1>{t("title")}</h1>;
};
