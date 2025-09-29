import type { ErrorResponseCodeType } from '@Front/types/api.types';

export type OAuthProvidersResponseType = {
  discord: string;
  github: string;
  google: string;
};

export type OAuthProvidersErrorCodeType = ErrorResponseCodeType<'PROVIDER_CONNECTION_FAILED'>;
