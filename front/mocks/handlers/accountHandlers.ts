import { account, accountError } from '@Mocks/fixtures/accountFixtures';
import { delay, http, HttpResponse } from 'msw';

export const postAccount201 = http.post(`http://localhost:3005/v1/account`, async () => {
  await delay();

  return HttpResponse.json(account, { status: 201 });
});

export const postAccount400 = http.post(`http://localhost:3005/v1/account`, async () => {
  await delay();

  return HttpResponse.json(accountError, { status: 400 });
});