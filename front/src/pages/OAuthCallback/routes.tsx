import type { RouteObject } from 'react-router';
import { routeAuthRegistry } from '@Front/utils/routeAuthRegistry';
import { OAuthCallback } from './OAuthCallback';

// Register this route as not requiring authentication
routeAuthRegistry.registerUnauthenticatedPath('oauth/callback');

export const oauthCallbackRoutes: RouteObject = {
  path: 'oauth/callback',
  element: <OAuthCallback />,
};
