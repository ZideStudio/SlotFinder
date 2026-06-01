import { Spinner } from "@Front/ui/atoms/Spinner/Spinner";

import "./Loader.scss";

export const Loader = () => (
  <div className="loader">
    <Spinner className="loader__spinner" />
  </div>
);
