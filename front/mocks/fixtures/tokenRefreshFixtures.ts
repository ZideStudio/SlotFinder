import type { ErrorResponseType } from "@Front/types/api.types";
import type { TokenRefreshErrorCodeType } from "@Front/types/Authentication/tokenRefresh/tokenRefresh.types";
import { SERVER_ERROR } from "@Front/utils/constants/api";

export const postTokenRefresh200Fixture = null;

export const postTokenRefresh500Fixture: ErrorResponseType<TokenRefreshErrorCodeType> =
  { code: SERVER_ERROR };
