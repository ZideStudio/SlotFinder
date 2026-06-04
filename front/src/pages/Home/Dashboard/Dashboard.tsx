import { Header } from "@Front/components/Layout/Header/Header";
import { useTranslation } from "react-i18next";

import "./Dashboard.scss";
import { Button } from "@Front/ui/molecules/Button/Button";
import { useLoader } from "@Front/hooks/useLoader";

export const Dashboard = () => {
  const { t } = useTranslation("dashboard");
  const { showLoader, hideLoader } = useLoader();

  const displayLoader = () => {
    showLoader();
    setTimeout(() => {
      hideLoader();
    }, 2000);
  };

  return (
    <div className="dashboard">
      <Header ignoreRouteHideHeader className="dashboard__header" />
      <div className="dashboard__content">
        <h1>{t("title")}</h1>

        <Button variant="primary" onClick={displayLoader}>
          Click Me
        </Button>
      </div>
    </div>
  );
};
