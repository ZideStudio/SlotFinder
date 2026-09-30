import type { RouteObject } from "react-router";
import { WhoAreYou } from "./WhoAreYou";

export const whoAreYouRoutes: RouteObject = {
  path: "who-are-you",
  element: <WhoAreYou />,
  handle: {
    hideHeader: true,
    mustBeAuthenticate: false,
  },
};
