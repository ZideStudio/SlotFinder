import type { PatchAccountErrorCodeType } from "@Front/api/account/patchAccount/types";
import { type HelpersApiError } from "@Front/api/generated/slotFinderAPI.schemas";
import { SERVER_ERROR } from "@Front/utils/constants/api";
import {
  getAccountMe200Fixture,
  patchAccount200Fixture,
  patchAccount400Fixture,
  postAccount201Fixture,
  postAccount400Fixture,
} from "@Mocks/fixtures/accountFixtures";
import { delay, http, HttpResponse } from "msw";

export const postAccount201 = http.post(
  `${import.meta.env.FRONT_BACKEND_URL}/v1/account`,
  async () => {
    await delay();

    return HttpResponse.json(postAccount201Fixture, { status: 201 });
  },
);

export const postAccount400 = http.post(
  `${import.meta.env.FRONT_BACKEND_URL}/v1/account`,
  async () => {
    await delay();

    return HttpResponse.json(postAccount400Fixture, { status: 400 });
  },
);

export const patchAccount200 = http.patch(
  `${import.meta.env.FRONT_BACKEND_URL}/v1/account`,
  async () => {
    await delay();

    return HttpResponse.json(patchAccount200Fixture, { status: 200 });
  },
);

export const patchAccount400 = (errorCode: PatchAccountErrorCodeType) =>
  http.patch(`${import.meta.env.FRONT_BACKEND_URL}/v1/account`, async () => {
    await delay();

    return HttpResponse.json(patchAccount400Fixture(errorCode), {
      status: 400,
    });
  });

export const patchAvatarAccount200 = http.patch(
  `${import.meta.env.FRONT_BACKEND_URL}/v1/account/avatar`,
  async () => {
    await delay();

    return HttpResponse.json(patchAccount200Fixture, { status: 200 });
  },
);

export const patchAvatarAccount400 = http.patch(
  `${import.meta.env.FRONT_BACKEND_URL}/v1/account/avatar`,
  async () => {
    await delay();

    return HttpResponse.json<HelpersApiError>(
      { code: SERVER_ERROR },
      { status: 400 },
    );
  },
);

export const getAccountMe200 = http.get(
  `${import.meta.env.FRONT_BACKEND_URL}/v1/account/me`,
  async () => {
    await delay();

    return HttpResponse.json(getAccountMe200Fixture, { status: 200 });
  },
);

export const getAccountAvatar200 = http.get(
  `${import.meta.env.FRONT_BACKEND_URL}/v1/account/:id/avatar`,
  async () => {
    await delay();

    return new HttpResponse(new Uint8Array(0), {
      status: 200,
      headers: { "Content-Type": "image/png" },
    });
  },
);
