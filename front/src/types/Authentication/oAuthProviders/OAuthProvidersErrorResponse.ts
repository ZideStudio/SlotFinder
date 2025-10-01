import { ErrorResponse } from '@Front/types/ErrorResponse';
import type { OAuthProvidersErrorCodeType } from './oAuthProviders.types';

export class OAuthProvidersErrorResponse extends ErrorResponse<OAuthProvidersErrorCodeType> {}
