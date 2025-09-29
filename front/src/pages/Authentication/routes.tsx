import type { RouteObject } from 'react-router-dom';
import { Authentication } from './Authentication';
import { signUpRoutes } from './SignUp';

export const authenticationRoutes: RouteObject = {
  element: <Authentication />,
  children: [signUpRoutes],
};
