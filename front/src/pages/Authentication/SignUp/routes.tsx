import type { RouteObject } from 'react-router';
import { SignUp } from './SignUp';

export const signUpRoutes: RouteObject = {
  path: 'sign-up',
  element: <SignUp />,
  handle: {
    mustBeAuthenticate: false,
  },
};
