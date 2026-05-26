import type { RouteObject } from 'react-router';
import { Loader } from './Loader';


export const loaderRoutes: RouteObject = {
  path: '/loader',
  index: true,
  element: <Loader />,
  handle: {
    hideHeader: true,
  },
};
