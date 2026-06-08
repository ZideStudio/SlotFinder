import { Header } from "@Front/components/Layout/Header/Header";

import "./Dashboard.scss";
import { Heading } from "@Front/ui/atoms/Heading/Heading";
import { Button } from "@Front/ui/molecules/Button/Button";

export const Dashboard = () => (
  <div className="dashboard">
    <Header ignoreRouteHideHeader className="dashboard__header" />
    <div className="dashboard__content">
      <div className="dashboard__content--header">
        <Heading level={1}>My events</Heading>

        <Button className="dashboard__button">Create an event</Button>
      </div>

      <div className="dashboard__content--no-events">No events here</div>
    </div>
  </div>
);
