import { signUpRoutes } from '@Front/pages/Authentication/SignUp';
import { errorRoutes } from '@Front/pages/Error';
import { oauthCallbackRoutes } from '@Front/pages/OAuthCallback';

type AppRoute = {
  home: () => string;
  signUp: () => string;
  oAuthCallback: (params?: { error?: string; returnUrl?: string }) => string;
  error: () => string;
};

export const appRoutes: AppRoute = {
  home: () => '/',
  signUp: () => `/${signUpRoutes.path}`,
  oAuthCallback: ({ error, returnUrl } = {}) => {
    let route = `/${oauthCallbackRoutes.path}`;

    const queryParams = new URLSearchParams();

    if (error) {
      queryParams.append('error', error);
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
