import type { RouteObject } from 'react-router-dom';
import { Dashboard } from './Dashboard';

export const dashboardRoutes: RouteObject = {
  path: 'dashboard',
  element: <Dashboard />,
  handle: {
    mustBeAuthenticate: true,
  },
};
