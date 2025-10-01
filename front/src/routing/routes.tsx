import { Layout } from '@Front/components/Layout';
import { authenticationRoutes } from '@Front/pages/Authentication';
import { dashboardRoutes } from '@Front/pages/Dashboard';
import { errorRoutes } from '@Front/pages/Error';
import { homeRoutes } from '@Front/pages/Home';
import { oauthCallbackRoutes } from '@Front/pages/OAuthCallback';
import type { RouteObject } from 'react-router-dom';

export const routeObject: RouteObject[] = [
  {
    path: '/',
    element: <Layout />,
    children: [homeRoutes, dashboardRoutes, authenticationRoutes, oauthCallbackRoutes, errorRoutes],
  },
];
