import { signUpRoutes } from "@Front/pages/Authentication/SignUp";
import { forgotPasswordRoutes } from "@Front/pages/ForgotPassword/routes";
import { errorRoutes } from "@Front/pages/Error";
import { oauthCallbackRoutes } from "@Front/pages/OAuthCallback";
import { whoAreYouRoutes } from "@Front/pages/WhoAreYou/routes";

export const appRoutes = {
  home: () => "/",
  signUp: () => `/${signUpRoutes.path}`,
  forgotPassword: () => `/${forgotPasswordRoutes.path}`,
  oAuthCallback: ({
    error,
    returnUrl,
  }: { error?: string; returnUrl?: string } = {}) => {
    let route = `/${oauthCallbackRoutes.path}`;

    const queryParams = new URLSearchParams();

    if (error) {
      queryParams.append("error", error);
    }
    if (returnUrl) {
      queryParams.append("returnUrl", returnUrl);
    }
    if (queryParams.toString()) {
      route += `?${queryParams.toString()}`;
    }
    return route;
  },
  whoAreYou: () => `/${whoAreYouRoutes.path}`,
  error: () => `/${errorRoutes.path}`,
};
