import type { ErrorResponseType } from "@Front/types/api.types";
import type {
  AuthStatusErrorCodeType,
  AuthStatusResponseType,
} from "@Front/types/Authentication/authStatus/authStatus.types";

export const authStatus200Fixture: AuthStatusResponseType = null;

export const authStatus401Fixture: ErrorResponseType<AuthStatusErrorCodeType> =
  { code: "NOT_AUTHENTICATED" };

export const authStatus403Fixture: ErrorResponseType<AuthStatusErrorCodeType> =
  { code: "TERMS_NOT_ACCEPTED" };
