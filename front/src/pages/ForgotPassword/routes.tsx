import type { RouteObject } from "react-router";
import { ForgotPassword } from "./ForgotPassword";

export const forgotPasswordRoutes: RouteObject & { path: string } = {
  path: "forgot-password",
  element: <ForgotPassword />,
  handle: {
    mustBeAuthenticate: false,
  },
};
