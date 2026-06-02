import type { ErrorResponseType } from "@Front/types/api.types";
import type {
  AuthStatusErrorCodeType,
  AuthStatusResponseType,
} from "@Front/types/Authentication/authStatus/authStatus.types";

export const getAuthStatus200Fixture: AuthStatusResponseType = null;

export const getAuthStatus401Fixture: ErrorResponseType<AuthStatusErrorCodeType> =
  { code: "NOT_AUTHENTICATED" };

export const getAuthStatus403Fixture: ErrorResponseType<AuthStatusErrorCodeType> =
  { code: "TERMS_NOT_ACCEPTED" };

export const getAuthStatus498Fixture: ErrorResponseType<AuthStatusErrorCodeType> =
  { code: "TOKEN_EXPIRED" };
