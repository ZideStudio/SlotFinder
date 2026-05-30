import {
  getAuthStatus200Fixture,
  getAuthStatus401Fixture,
  getAuthStatus403Fixture,
  getAuthStatus498Fixture,
} from "@Mocks/fixtures/authStatusFixtures";
import { delay, http, HttpResponse } from "msw";

export const getAuthStatus200 = http.get(
  `${import.meta.env.FRONT_BACKEND_URL}/v1/auth/status`,
  async () => {
    await delay();

    return HttpResponse.json(getAuthStatus200Fixture, { status: 200 });
  },
);

export const getAuthStatus401 = http.get(
  `${import.meta.env.FRONT_BACKEND_URL}/v1/auth/status`,
  async () => {
    await delay();

    return HttpResponse.json(getAuthStatus401Fixture, { status: 401 });
  },
);

export const getAuthStatus403 = http.get(
  `${import.meta.env.FRONT_BACKEND_URL}/v1/auth/status`,
  async () => {
    await delay();

    return HttpResponse.json(getAuthStatus403Fixture, { status: 403 });
  },
);

export const getAuthStatus498 = http.get(
  `${import.meta.env.FRONT_BACKEND_URL}/v1/auth/status`,
  async () => {
    await delay();

    return HttpResponse.json(getAuthStatus498Fixture, { status: 498 });
  },
);
