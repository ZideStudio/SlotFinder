import { accountErrorFixture, accountFixture } from '@Mocks/fixtures/accountFixtures';
import { delay, http, HttpResponse } from 'msw';

export const postAccount201 = http.post(`${import.meta.env.FRONT_BACKEND_URL}/v1/account`, async () => {
  await delay();

  return HttpResponse.json(accountFixture, { status: 201 });
});

export const postAccount400 = http.post(`${import.meta.env.FRONT_BACKEND_URL}/v1/account`, async () => {
  await delay();

  return HttpResponse.json(accountErrorFixture, { status: 400 });
});
