import type { ErrorResponseType } from "@Front/types/api.types";
import type { AuthStatusErrorCodeType, AuthStatusResponseType } from "@Front/types/Authentication/authStatus/authStatus.types";

export const authStatusFixture: AuthStatusResponseType = null;

export const authStatusErrorFixture: ErrorResponseType<AuthStatusErrorCodeType> = {
  code: 'NOT_AUTHENTICATED',
};
