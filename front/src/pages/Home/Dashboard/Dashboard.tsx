import { Heading } from "@Front/ui/atoms/Heading/Heading";
import { Button } from "@Front/ui/molecules/Button/Button";

import "./Dashboard.scss";

export const Dashboard = () => (
  <div className="dashboard">
    <div className="dashboard__header">
      <Heading level={1}>My events</Heading>

      <Button className="dashboard__button">Create an event</Button>
    </div>

    <div className="dashboard__no-events">No events here</div>
  </div>
);
