import { ErrorResponse } from '@Front/types/ErrorResponse';
import type { SignUpErrorCodeType } from './signUp.types';

export class SignUpErrorResponse extends ErrorResponse<SignUpErrorCodeType> {}
