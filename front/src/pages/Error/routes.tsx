import type { RouteObject } from 'react-router';
import { ErrorPage } from './Error';

export const errorRoutes: RouteObject = {
  path: 'error',
  element: <ErrorPage />,
};
