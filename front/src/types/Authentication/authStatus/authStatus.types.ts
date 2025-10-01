import type { ErrorResponseCodeType } from '@Front/types/api.types';

export type AuthStatusResponseType = null;

export type AuthStatusErrorCodeType = ErrorResponseCodeType<'NOT_AUTHENTICATED' | 'TOKEN_INVALID' | 'TOKEN_EXPIRED'>;
