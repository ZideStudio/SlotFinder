import type { RouteObject } from 'react-router-dom';
import { SignUp } from './SignUp';

export const signUpRoutes: RouteObject = {
  path: 'sign-up',
  element: <SignUp />,
  handle: {
    mustBeAuthenticate: false,
  },
};
