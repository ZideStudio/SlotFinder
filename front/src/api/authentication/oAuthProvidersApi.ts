import { appRoutes } from '@Front/routing/appRoutes';
import type {
  OAuthProvidersErrorCodeType,
  OAuthProvidersResponseType,
} from '@Front/types/Authentication/oAuthProviders/oAuthProviders.types';
import { OAuthProvidersErrorResponse } from '@Front/types/Authentication/oAuthProviders/OAuthProvidersErrorResponse';
import { METHODS } from '../constant';
import { fetchApi } from '../fetchApi';

export const oAuthProvidersApi = async () =>
  await fetchApi<OAuthProvidersResponseType, OAuthProvidersErrorCodeType>({
    path: `${import.meta.env.FRONT_BACKEND_URL}/v1/auth/providers/url?redirectUrl=${encodeURIComponent(import.meta.env.FRONT_DOMAIN + appRoutes.oAuthCallback())}`,
    method: METHODS.get,
    CustomErrorResponse: OAuthProvidersErrorResponse,
  });
