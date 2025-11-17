import { AuthenticationProtection } from '@Front/components/AuthenticationProtection/AuthenticationProtection';
import { Layout } from '@Front/components/Layout';
import { authenticationRoutes } from '@Front/pages/Authentication';
import { errorRoutes } from '@Front/pages/Error';
import { homeRoutes } from '@Front/pages/Home';
import { oauthCallbackRoutes } from '@Front/pages/OAuthCallback';
import type { RouteObject } from 'react-router';

export const routeObject: RouteObject[] = [
  {
    path: '/',
    element: (
      <AuthenticationProtection>
        <Layout />
      </AuthenticationProtection>
    ),
    children: [homeRoutes, authenticationRoutes, oauthCallbackRoutes, errorRoutes],
  },
];
