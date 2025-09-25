import { Layout } from '@Front/components/Layout';
import { homeRoutes } from '@Front/pages/Home';
import { signUpRoutes } from '@Front/pages/SignUp';
import type { RouteObject } from 'react-router';

export const routeObject: RouteObject[] = [
  {
    path: '/',
    element: <Layout />,
    children: [homeRoutes, signUpRoutes],
  },
];
