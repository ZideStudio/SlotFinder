import type { ErrorResponseCodeType } from "@Front/types/api.types";

export type AuthStatusResponseType = null;

export type AuthStatusErrorCodeType = ErrorResponseCodeType<
  "NOT_AUTHENTICATED" | "TERMS_NOT_ACCEPTED" | "TOKEN_INVALID" | "TOKEN_EXPIRED"
>;
