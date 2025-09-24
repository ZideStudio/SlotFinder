import { homeRoutes } from '@Front/pages/Home';
import { signUpRoutes } from '@Front/pages/SignUp';
import type { RouteObject } from 'react-router';

export const routeObject: RouteObject[] = [
  {
    path: '/',
    children: [homeRoutes, signUpRoutes],
  },
];
