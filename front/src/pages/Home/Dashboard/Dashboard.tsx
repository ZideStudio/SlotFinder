import { Header } from "@Front/components/Layout/Header/Header";
import { Heading } from "@Front/ui/atoms/Heading/Heading";
import { Icon } from "@Front/ui/atoms/Icon/Icon";
import { Button } from "@Front/ui/molecules/Button/Button";
import BoatIcon from "@material-symbols/svg-300/outlined/directions_boat.svg";

import { useTranslation } from "react-i18next";
import "./Dashboard.scss";

export const Dashboard = () => {
  const { t } = useTranslation("dashboard");
  return (
    <div className="dashboard">
      <Header ignoreRouteHideHeader className="dashboard__header" />
      <div className="dashboard__content">
        <div className="dashboard__content--header">
          <Heading level={1}>{t("title")}</Heading>

          <Button className="dashboard__content--header-buttons">
            {t("Create an event")}
          </Button>
        </div>
        <section className="dashboard__content--no-events">
          <div>{t("No events")}</div>
          <Icon
            className="dashboard__content--no-events-icon"
            icon={BoatIcon}
          />
        </section>
      </div>
    </div>
  );
};
