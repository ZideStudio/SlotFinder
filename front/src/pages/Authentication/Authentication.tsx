import { Outlet } from "react-router";

import "./Authentication.scss";
import { Heading } from "@Front/ui/atoms/Heading/Heading";

export const Authentication = () => (
  <section className="authentication subgrid">
    <Heading level={1} className="authentication__title">
      Authentication Page
    </Heading>
    <div className="authentication__content">
      <Outlet />
    </div>
  </section>
);
