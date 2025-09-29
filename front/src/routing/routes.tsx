import { Layout } from '@Front/components/Layout';
import { authenticationRoutes } from '@Front/pages/Authentication';
import { errorRoutes } from '@Front/pages/Error';
import { homeRoutes } from '@Front/pages/Home';
import { oauthCallbackRoutes } from '@Front/pages/OAuthCallback';
import type { RouteObject } from 'react-router-dom';

export const routeObject: RouteObject[] = [
  {
    path: '/',
    element: <Layout />,
    children: [homeRoutes, authenticationRoutes, oauthCallbackRoutes, errorRoutes],
  },
];
