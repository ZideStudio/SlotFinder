import {
  authStatus200Fixture,
  authStatus401Fixture,
  authStatus403Fixture,
} from "@Mocks/fixtures/authStatusFixtures";
import { delay, http, HttpResponse } from "msw";

export const getAuthStatus200 = http.get(
  `${import.meta.env.FRONT_BACKEND_URL}/v1/auth/status`,
  async () => {
    await delay();

    return HttpResponse.json(authStatus200Fixture, { status: 200 });
  },
);

export const getAuthStatus401 = http.get(
  `${import.meta.env.FRONT_BACKEND_URL}/v1/auth/status`,
  async () => {
    await delay();

    return HttpResponse.json(authStatus401Fixture, { status: 401 });
  },
);

export const getAuthStatus403 = http.get(
  `${import.meta.env.FRONT_BACKEND_URL}/v1/auth/status`,
  async () => {
    await delay();

    return HttpResponse.json(authStatus403Fixture, { status: 403 });
  },
);
