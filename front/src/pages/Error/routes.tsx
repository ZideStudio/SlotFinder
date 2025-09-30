import type { RouteObject } from 'react-router-dom';
import { ErrorPage } from './Error';

export const errorRoutes: RouteObject = {
  path: 'error',
  element: <ErrorPage />,
};
