import {
  postTokenRefresh200Fixture,
  postTokenRefresh500Fixture,
} from "@Mocks/fixtures/tokenRefreshFixtures";
import { delay, http, HttpResponse } from "msw";

export const postTokenRefresh200 = http.post(
  `${import.meta.env.FRONT_BACKEND_URL}/v1/auth/refresh`,
  async () => {
    await delay();

    return HttpResponse.json(postTokenRefresh200Fixture, { status: 200 });
  },
);

export const postTokenRefresh500 = http.post(
  `${import.meta.env.FRONT_BACKEND_URL}/v1/auth/refresh`,
  async () => {
    await delay();

    return HttpResponse.json(postTokenRefresh500Fixture, { status: 500 });
  },
);

export const postTokenRefreshNetworkError = http.post(
  `${import.meta.env.FRONT_BACKEND_URL}/v1/auth/refresh`,
  async () => {
    await delay();

    return HttpResponse.error();
  },
);
