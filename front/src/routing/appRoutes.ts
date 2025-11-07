import { signUpRoutes } from '@Front/pages/Authentication/SignUp';
import { dashboardRoutes } from '@Front/pages/Dashboard/routes';
import { errorRoutes } from '@Front/pages/Error';
import { oauthCallbackRoutes } from '@Front/pages/OAuthCallback';

type AppRoute = {
  home: () => string;
  dashboard: () => string;
  signUp: () => string;
  oAuthCallback: (params?: { message?: string; returnUrl?: string }) => string;
  error: () => string;
};

export const appRoutes: AppRoute = {
  home: () => '/',
  dashboard: () => `/${dashboardRoutes.path}`,
  signUp: () => `/${signUpRoutes.path}`,
  oAuthCallback: ({ message, returnUrl } = {}) => {
    let route = `/${oauthCallbackRoutes.path}`;

    const queryParams = new URLSearchParams();

    if (message) {
      queryParams.append('message', message);
    }
    if (returnUrl) {
      queryParams.append('returnUrl', returnUrl);
    }
    if (queryParams.toString()) {
      route += `?${queryParams.toString()}`;
    }
    return route;
  },
  error: () => `/${errorRoutes.path}`,
};
