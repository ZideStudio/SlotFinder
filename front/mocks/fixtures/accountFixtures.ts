import type { HelpersApiError } from "@Front/api/generated/slotFinderAPI.schemas";
import type { SignUpResponseType } from "@Front/types/Authentication/signUp/signUp.types";

// The OpenAPI spec documents the POST /v1/account response as void, but
// The real API returns a token object. The shape below matches the real API.
export const postAccount201Fixture: SignUpResponseType = {
  access_token: "1234567890abcdef",
  createdAt: "2024-01-01T00:00:00.000Z",
  email: "test@example.com",
  id: "123456",
  providers: null,
  userName: "test_user",
};

export const postAccount400Fixture: HelpersApiError = {
  code: "USERNAME_ALREADY_TAKEN",
};
