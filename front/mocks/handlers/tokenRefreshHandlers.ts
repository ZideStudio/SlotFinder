import { delay, http, HttpResponse } from 'msw';

export const postTokenRefresh200 = (delayMs?: number) =>
  http.post(`${import.meta.env.FRONT_BACKEND_URL}/v1/auth/refresh`, async () => {
    await delay(delayMs);

    return HttpResponse.json({}, { status: 200 });
  });

export const postTokenRefresh400 = http.post(`${import.meta.env.FRONT_BACKEND_URL}/v1/auth/refresh`, async () => {
  await delay();

  return HttpResponse.json({ error: 'Refresh token expired' }, { status: 400 });
});

export const postTokenRefreshNetworkError = http.post(`${import.meta.env.FRONT_BACKEND_URL}/v1/auth/refresh`, async () => {
  await delay();

  return HttpResponse.error();
});
