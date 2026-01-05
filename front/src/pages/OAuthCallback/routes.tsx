import type { RouteObject } from 'react-router';
import { OAuthCallback } from './OAuthCallback';

export const oauthCallbackRoutes: RouteObject = {
  path: 'oauth/callback',
  element: <OAuthCallback />,
};
