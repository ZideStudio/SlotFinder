import { Layout } from '@Front/components/Layout';
import { authenticationRoutes } from '@Front/pages/Authentication';
import { homeRoutes } from '@Front/pages/Home';
import type { RouteObject } from 'react-router';

export const routeObject: RouteObject[] = [
  {
    path: '/',
    element: <Layout />,
    children: [homeRoutes, authenticationRoutes],
  },
];
