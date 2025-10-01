import type { ErrorResponseType } from '@Front/types/api.types';
import type {
  OAuthProvidersErrorCodeType,
  OAuthProvidersResponseType,
} from '@Front/types/Authentication/oAuthProviders/oAuthProviders.types';

export const oAuthProvidersFixture: OAuthProvidersResponseType = {
  discord: 'https://discord.com/oauth/',
  github: 'https://github.com/oauth/',
  google: 'https://google.com/oauth/',
};

export const oAuthProvidersErrorFixture: ErrorResponseType<OAuthProvidersErrorCodeType> = {
  code: 'PROVIDER_CONNECTION_FAILED',
};
