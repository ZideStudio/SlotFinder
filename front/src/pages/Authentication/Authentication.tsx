import { Outlet } from "react-router";
import { OAuth } from "./OAuth/OAuth";

import "./Authentication.scss";

export const Authentication = () => (
  <section className="authentication subgrid">
    <h1 className="authentication__title">Authentication Page</h1>
    <div className="authentication__content">
      <Outlet />
    </div>
    <OAuth />
  </section>
);
