import { signUpRoutes } from '@Front/pages/Authentication/SignUp';
import { dashboardRoutes } from '@Front/pages/Dashboard/routes';
import { errorRoutes } from '@Front/pages/Error';
import { oauthCallbackRoutes } from '@Front/pages/OAuthCallback';

type AppRoute = {
  home: () => string;
  dashboard: () => string;
  signUp: () => string;
  oAuthCallback: (message?: string) => string;
  error: () => string;
};

export const appRoutes: AppRoute = {
  home: () => '/',
  dashboard: () => `/${dashboardRoutes.path}`,
  signUp: () => `/${signUpRoutes.path}`,
  oAuthCallback: message => {
    let route = `/${oauthCallbackRoutes.path}`;
    if (message) {
      route += `?message=${encodeURIComponent(message)}`;
    }
    return route;
  },
  error: () => `/${errorRoutes.path}`,
};
