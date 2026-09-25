import type { PatchAccountErrorCodeType } from "@Front/api/account/patchAccount/types";
import type {
  AccountAccountResponseDto,
  HelpersApiError,
} from "@Front/api/generated/slotFinderAPI.schemas";
import { TERMS_VERSION } from "@Front/utils/constants/terms";

export const postAccount201Fixture = undefined;

export const postAccount400Fixture: HelpersApiError = {
  code: "USERNAME_ALREADY_TAKEN",
};

export const patchAccount200Fixture: AccountAccountResponseDto = {
  username: "test",
  email: "test@example.com",
  avatarUrl: "/api/v1/account/123456789/avatar",
  color: "#ff0000",
  language: "en",
  providers: [
    {
      provider: "github",
    },
  ],
  timeZone: "Europe/Paris",
};

export const patchAccount400Fixture = (
  errorCode: PatchAccountErrorCodeType,
): HelpersApiError => ({
  code: errorCode,
});

export const getAccountMe200Fixture: AccountAccountResponseDto = {
  username: "test_user",
  email: "test@example.com",
  avatarUrl: "/api/v1/account/123456789/avatar",
  color: "#ff0000",
  language: "en",
  providers: [
    {
      provider: "github",
    },
  ],
  timeZone: "Europe/Paris",
  termsAccepted: true,
  termsVersion: TERMS_VERSION,
};

export const getAccountMeWithoutTerms200Fixture: AccountAccountResponseDto = {
  username: "test_user",
  email: "test@example.com",
  avatarUrl: "/api/v1/account/123456789/avatar",
  color: "#ff0000",
  language: "en",
  providers: [
    {
      provider: "github",
    },
  ],
  timeZone: "Europe/Paris",
  termsAccepted: false,
  termsVersion: undefined,
};
