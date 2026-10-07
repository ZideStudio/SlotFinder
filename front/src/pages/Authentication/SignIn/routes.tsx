import type { RouteObject } from "react-router";
import { SignIn } from "./SignIn";

export const signInRoutes: RouteObject = {
  path: "sign-in",
  element: <SignIn />,
  handle: {
    mustBeAuthenticate: false,
  },
};
