import type { HelpersApiError } from "@Front/api/generated/slotFinderAPI.schemas";

export const getAuthStatus200Fixture = null;

export const getAuthStatus401Fixture: HelpersApiError = {
  code: "NOT_AUTHENTICATED",
};

export const getAuthStatus403Fixture: HelpersApiError = {
  code: "TERMS_NOT_ACCEPTED",
};

export const getAuthStatus498Fixture: HelpersApiError = {
  code: "TOKEN_EXPIRED",
};
