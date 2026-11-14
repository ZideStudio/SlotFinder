import type { RouteObject } from "react-router";
import { Authentication } from "./Authentication";
import { forgotPasswordRoutes } from "../ForgotPassword/routes";
import { signUpRoutes } from "./SignUp";

export const authenticationRoutes: RouteObject = {
  element: <Authentication />,
  children: [signUpRoutes, forgotPasswordRoutes],
  handle: {
    hideHeader: true,
  },
};
