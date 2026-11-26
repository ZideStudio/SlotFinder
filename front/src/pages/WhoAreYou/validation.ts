import { mbToBytes } from "@Front/utils/helpers/fileSizeConversions";
import type { TFunction } from "i18next";
import { z } from "zod";
import {
  AVATAR_FILE_TYPES,
  AVATAR_MAX_SIZE_MB,
  COLOR_REGEX,
  USERNAME_MAX_LENGTH,
  USERNAME_MIN_LENGTH,
} from "./constants";
import type { WhoAreYouFormData } from "./types";

export const getSchema = (translate: TFunction<"whoAreYou">) =>
  z.object({
    avatar: z
      .custom<FileList>(
        (value) =>
          typeof FileList !== "undefined" &&
          value instanceof FileList &&
          value.length > 0,
        { message: translate("avatarRequired") },
      )
      .refine(
        (value) =>
          !(typeof FileList !== "undefined" && value instanceof FileList) ||
          value.length === 0 ||
          value[0].size <= mbToBytes(AVATAR_MAX_SIZE_MB),
        {
          message: translate("avatarFileSizeError", {
            maxSize: AVATAR_MAX_SIZE_MB,
          }),
        },
      )
      .refine(
        (value) =>
          !(typeof FileList !== "undefined" && value instanceof FileList) ||
          value.length === 0 ||
          AVATAR_FILE_TYPES.includes(value[0].type),
        { message: translate("avatarFileTypeError") },
      ),
    username: z
      .string()
      .min(1, translate("requiredUsername"))
      .min(
        USERNAME_MIN_LENGTH,
        translate("minLengthUsername", { min: USERNAME_MIN_LENGTH }),
      )
      .max(
        USERNAME_MAX_LENGTH,
        translate("maxLengthUsername", { max: USERNAME_MAX_LENGTH }),
      ),
    color: z
      .string()
      .min(1, translate("colorRequired"))
      .regex(COLOR_REGEX, translate("colorInvalid")),
    termsAccepted: z.custom<boolean>((value) => value === true, {
      message: translate("termsAcceptedError"),
    }),
  }) satisfies z.ZodType<WhoAreYouFormData>;
