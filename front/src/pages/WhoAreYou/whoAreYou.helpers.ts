import type { AccountAccountUpdateDto } from "@Front/api/generated/slotFinderAPI.schemas";
import type { WhoAreYouFormData } from "./types";

export const buildAccountUpdateDTO = (
  data: WhoAreYouFormData,
): AccountAccountUpdateDto => ({
  username: data.username,
  color: data.color,
  termsAccepted: data.termsAccepted,
  timeZone: Temporal.Now.timeZoneId(),
});
