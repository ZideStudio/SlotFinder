import type { RouteObject } from 'react-router';
import { routeAuthRegistry } from '@Front/utils/routeAuthRegistry';
import { SignUp } from './SignUp';

// Register this route as not requiring authentication
routeAuthRegistry.registerUnauthenticatedPath('sign-up');

export const signUpRoutes: RouteObject = {
  path: 'sign-up',
  element: <SignUp />,
  handle: {
    mustBeAuthenticate: false,
  },
};
