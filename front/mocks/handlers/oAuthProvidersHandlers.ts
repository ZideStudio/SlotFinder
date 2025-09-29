import { oAuthProvidersErrorFixture, oAuthProvidersFixture } from '@Mocks/fixtures/oAuthProvidersFixtures';
import { delay, http, HttpResponse } from 'msw';

export const getOAuthProviders200 = http.get(`${import.meta.env.FRONT_BACKEND_URL}/v1/auth/providers/url`, async () => {
  await delay();

  return HttpResponse.json(oAuthProvidersFixture, { status: 200 });
});

export const getOAuthProviders400 = http.get(`${import.meta.env.FRONT_BACKEND_URL}/v1/auth/providers/url`, async () => {
  await delay();

  return HttpResponse.json(oAuthProvidersErrorFixture, { status: 400 });
});
