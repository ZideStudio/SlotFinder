import { postV1Account } from "@Front/api/generated/account/account";
import {
  AccountAccountCreateDtoLanguage,
  type AccountAccountCreateDto,
} from "@Front/api/generated/slotFinderAPI.schemas";
import type { SignUpRequestBodyType } from "@Front/types/Authentication/signUp/signUp.types";
import { TERMS_VERSION } from "@Front/utils/constants/terms";
import { Temporal } from "@js-temporal/polyfill";

/**
 * Narrows a generic `string` to `AccountAccountCreateDtoLanguage` ('en' | 'fr').
 * Returns `true` only when the value is one of the enum values at runtime.
 */
const isValidLanguage = (
  lang: string,
): lang is AccountAccountCreateDtoLanguage => {
  const values: readonly string[] = Object.values(
    AccountAccountCreateDtoLanguage,
  );
  return values.includes(lang);
};

export const signUpApi = ({
  username,
  email,
  password,
  language,
}: SignUpRequestBodyType): Promise<void> => {
  const dto: AccountAccountCreateDto & { username: string } = {
    username,
    email,
    password,
    language: isValidLanguage(language)
      ? language
      : AccountAccountCreateDtoLanguage.en,
    termsAccepted: true,
    termsVersion: TERMS_VERSION,
    timeZone: Temporal.Now.timeZoneId(),
  };

  return postV1Account(dto);
};
