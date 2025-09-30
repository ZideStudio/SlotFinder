import { signUpRoutes } from '@Front/pages/Authentication/SignUp';
import { errorRoutes } from '@Front/pages/Error';
import { oauthCallbackRoutes } from '@Front/pages/OAuthCallback';

type AppRoute = {
  home: () => string;
  signUp: () => string;
  oAuthCallback: (message?: string) => string;
  error: () => string;
};

export const appRoutes: AppRoute = {
  home: () => '/',
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
