import { http, HttpResponse } from 'msw';
import { setupServer } from 'msw/node';

export const server = setupServer(
  http.post(`${import.meta.env.FRONT_BACKEND_URL}/v1/account`, () =>
    HttpResponse.json({
      access_token: '1234567890abcdef',
      createdAt: '2024-01-01T00:00:00.000Z',
      email: 'test@example.com',
      id: '123456',
      providers: null,
      userName: 'testuser',
    }),
  ),
);
