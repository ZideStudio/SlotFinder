import type { RouteObject } from 'react-router-dom';
import { Home } from './Home';

export const homeRoutes: RouteObject = {
  index: true,
  element: <Home />,
};
