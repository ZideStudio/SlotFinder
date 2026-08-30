import type { ErrorResponseCodeType } from "@Front/types/api.types";

export type AuthStatusErrorCodeType = ErrorResponseCodeType<
  | "NOT_AUTHENTICATED"
  | "USERNAME_MISSING"
  | "TERMS_NOT_ACCEPTED"
  | "TOKEN_INVALID"
  | "TOKEN_EXPIRED"
>;
