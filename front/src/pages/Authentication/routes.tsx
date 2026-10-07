import type { RouteObject } from "react-router";
import { Authentication } from "./Authentication";
import { signUpRoutes } from "./SignUp";
import { signInRoutes } from "./SignIn";

export const authenticationRoutes: RouteObject = {
  element: <Authentication />,
  children: [signUpRoutes, signInRoutes],
  handle: {
    hideHeader: true,
  },
};
