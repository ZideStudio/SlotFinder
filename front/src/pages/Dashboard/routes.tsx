import type { RouteObject } from 'react-router';
import { Dashboard } from './Dashboard';

export const dashboardRoutes: RouteObject = {
  path: 'dashboard',
  element: <Dashboard />,
  handle: {
    mustBeAuthenticate: true,
  },
};
