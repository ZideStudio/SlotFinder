import type { AccountAccountUpdateDto } from "@Front/api/generated/slotFinderAPI.schemas";
import { TERMS_VERSION } from "@Front/utils/constants/terms";
import type { WhoAreYouFormData } from "./types";

export const buildAccountUpdateDTO = (
  data: WhoAreYouFormData,
): AccountAccountUpdateDto => ({
  username: data.username,
  color: data.color,
  termsAccepted: data.termsAccepted,
  termsVersion: TERMS_VERSION,
  timeZone: Temporal.Now.timeZoneId(),
});
