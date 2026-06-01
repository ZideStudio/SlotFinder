import { Header } from "@Front/components/Layout/Header/Header";
import { useTranslation } from "react-i18next";

import "./Dashboard.scss";

export const Dashboard = () => {
  const { t } = useTranslation("dashboard");

  return (
    <div className="dashboard">
      <Header ignoreRouteHideHeader className="dashboard__header" />
      <div className="dashboard__content">
        <Heading level={1}>My events</Heading>

        <Button className="dashboard__button">Create an event</Button>
      </div>

      <div className="dashboard__no-events">No events here</div>
    </div>
  );
};
