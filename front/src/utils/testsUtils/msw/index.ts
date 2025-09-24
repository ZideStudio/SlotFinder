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
 * Méthodes permettant de mocker une route d'API en GET et POST depuis le fichier de test
 * @param base : BASE ROUTE avec comme valeur par défaut MOCK_API_URL.base
 * @param route : URI à renseigner
 * @param code : status code que l'API doit renvoyer
 * @param responseBody : corps de la réponse
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
