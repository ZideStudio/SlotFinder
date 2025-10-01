import { ErrorResponse } from '@Front/types/ErrorResponse';
import type { AuthStatusErrorCodeType } from './authStatus.types';

export class AuthStatusErrorResponse extends ErrorResponse<AuthStatusErrorCodeType> {}
