import { setupServer } from 'msw/node';
import { getAuthStatus200 } from './handlers/authStatusHandlers';
import { getOAuthProviders200 } from './handlers/oAuthProvidersHandlers';

export const server = setupServer(getAuthStatus200, getOAuthProviders200);
