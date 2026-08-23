import type { HelpersApiError } from "@Front/api/generated/slotFinderAPI.schemas";
import type { AuthStatusErrorCodeType } from "@Front/types/Authentication/authStatus/authStatus.types";

export const getAuthStatus200Fixture = null;

export const getAuthStatus401Fixture: HelpersApiError = {
  code: "NOT_AUTHENTICATED",
};

export const getAuthStatus403Fixture = (errorCode: AuthStatusErrorCodeType): HelpersApiError => ({
  code: errorCode,
});

export const getAuthStatus498Fixture: HelpersApiError = {
  code: "TOKEN_EXPIRED",
};
