import { authStatusErrorFixture, authStatusFixture } from '@Mocks/fixtures/authStatusFixtures';
import { delay, http, HttpResponse } from 'msw';

export const getAuthStatus200 = http.get(`${import.meta.env.FRONT_BACKEND_URL}/v1/auth/status`, async () => {
  await delay();

  return HttpResponse.json(authStatusFixture, { status: 200 });
});

export const getAuthStatus400 = http.get(`${import.meta.env.FRONT_BACKEND_URL}/v1/auth/status`, async () => {
  await delay();

  return HttpResponse.json(authStatusErrorFixture, { status: 400 });
});
