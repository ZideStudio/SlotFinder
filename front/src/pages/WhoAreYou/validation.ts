import { mbToBytes } from "@Front/utils/helpers/fileSizeConversions";
import type { TFunction } from "i18next";
import { boolean, mixed, object, string, type ObjectSchema } from "yup";
import {
  AVATAR_FILE_TYPES,
  AVATAR_MAX_SIZE_MB,
  COLOR_REGEX,
  USERNAME_MAX_LENGTH,
  USERNAME_MIN_LENGTH,
} from "./constants";
import type { WhoAreYouFormData } from "./types";

export const getSchema = (
  translate: TFunction,
): ObjectSchema<WhoAreYouFormData> =>
  object({
    avatar: mixed<FileList>()
      .required(translate("avatarRequired"))
      .test("fileList", translate("avatarRequired"), (value) => {
        if (!value || !(value instanceof FileList) || value.length === 0) {
          return false;
        }
        return true;
      })
      .test(
        "fileSize",
        translate("avatarFileSizeError", { maxSize: AVATAR_MAX_SIZE_MB }),
        (value) => {
          if (!value || !(value instanceof FileList) || value.length === 0) {
            return false;
          }
          return value[0].size <= mbToBytes(AVATAR_MAX_SIZE_MB);
        },
      )
      .test("fileType", translate("avatarFileTypeError"), (value) => {
        if (!value || !(value instanceof FileList) || value.length === 0) {
          return false;
        }
        return AVATAR_FILE_TYPES.includes(value[0].type);
      }),
    username: string()
      .required(translate("requiredUsername"))
      .min(
        USERNAME_MIN_LENGTH,
        translate("minLengthUsername", { min: USERNAME_MIN_LENGTH }),
      )
      .max(
        USERNAME_MAX_LENGTH,
        translate("maxLengthUsername", { max: USERNAME_MAX_LENGTH }),
      ),
    color: string()
      .required(translate("colorRequired"))
      .matches(COLOR_REGEX, translate("colorInvalid")),
    termsAccepted: boolean()
      .oneOf([true], translate("termsAcceptedError"))
      .required(translate("termsAcceptedError")),
  }).required();
