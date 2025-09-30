import { account, accountError } from '@Mocks/fixtures/accountFixtures';
import { delay, http, HttpResponse } from 'msw';

export const postAccount201 = http.post(`${import.meta.env.FRONT_BACKEND_URL}/v1/account`, async () => {
  await delay();

  return HttpResponse.json(account, { status: 201 });
});

export const postAccount400 = http.post(`${import.meta.env.FRONT_BACKEND_URL}/v1/account`, async () => {
  await delay();

  return HttpResponse.json(accountError, { status: 400 });
});
