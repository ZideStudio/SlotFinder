// oxlint-disable no-magic-numbers
import { server } from '@Mocks/server';
import { http, HttpResponse, type JsonBodyType } from 'msw';

type CommonResponseType = {
  code?: number;
  responseBody?: JsonBodyType;
};

export const commonResponse =
  ({ code = 200, responseBody }: CommonResponseType) =>
  () =>
    HttpResponse.json(responseBody, { status: code });

type ServerUseType = {
  base?: string;
  route?: string;
  code?: number;
  responseBody?: JsonBodyType;
};

/**
 * Methods to mock an API route for GET and POST from the test file
 * @param base: BASE ROUTE, default value is FRONT_BACKEND_URL
 * @param route: URI to specify
 * @param code: status code that the API should return
 * @param responseBody: response body
 */

export const serverUseGet = ({
  base = import.meta.env.FRONT_BACKEND_URL,
  route = '',
  code = 200,
  responseBody = {},
}: ServerUseType) => server.use(http.get(`${base}${route}`, commonResponse({ code, responseBody })));

export const serverUsePost = ({
  base = import.meta.env.FRONT_BACKEND_URL,
  route = '',
  code = 200,
  responseBody = {},
}: ServerUseType) => {
  server.use(http.post(`${base}${route}`, commonResponse({ code, responseBody })));
};
