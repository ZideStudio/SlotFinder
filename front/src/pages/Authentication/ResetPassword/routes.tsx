import type { RouteObject } from "react-router";
import { ResetPassword } from "./ResetPassword";

export const resetPasswordRoutes: RouteObject & { path: string } = {
  path: "reset-password",
  element: <ResetPassword />,
  handle: {
    mustBeAuthenticate: false,
    title: "resetPassword:title",
    hideOAuth: true,
  },
};
