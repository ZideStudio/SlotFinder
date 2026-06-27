import { Header } from "@Front/components/Layout/Header/Header";
import BoatIcon from "@material-symbols/svg-300/outlined/directions_boat.svg?react";
import { Heading } from "@Front/ui/atoms/Heading/Heading";
import { Button } from "@Front/ui/molecules/Button/Button";

import "./Dashboard.scss";
import { Icon } from "@Front/ui/atoms/Icon/Icon";

export const Dashboard = () => (
  <div className="dashboard">
    <Header ignoreRouteHideHeader className="dashboard__header" />
    <div className="dashboard__content">
      <div className="dashboard__content--header">
        <Heading level={1}>My events</Heading>

        <Button className="dashboard__content--header-buttons">
          Create an event
        </Button>
      </div>
      <div className="dashboard__content--no-events">
        <div>No events here</div>
        <Icon className="dashboard__content--no-events-icon" icon={BoatIcon} />
      </div>
    </div>
  </div>
);
