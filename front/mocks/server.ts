import { setupServer } from 'msw/node';
import { getAuthStatus200 } from './handlers/authStatusHandlers';

export const server = setupServer(getAuthStatus200);
