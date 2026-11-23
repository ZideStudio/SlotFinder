import {
  EMAIL_REGEX,
  PASSWORD_MIN_LENGTH,
  PASSWORD_REGEX,
} from "@Front/pages/Authentication/SignUp/constants";
import {
  AVATAR_FILE_TYPES,
  AVATAR_MAX_SIZE_MB,
  COLOR_REGEX,
  USERNAME_MAX_LENGTH,
  USERNAME_MIN_LENGTH,
} from "@Front/pages/WhoAreYou/constants";
import { type TFunction } from "i18next";
import { boolean, mixed, object, ref, string } from "yup";
import { mbToBytes } from "../helpers/fileSizeConversions";

export const getUserSchema = (
  translate: TFunction<
    "signUp" | "forgotPassword" | "whoAreYou"
  >,
) =>
  object({
    email: string()
      .required(translate("requiredEmail"))
      .max(320, translate("invalidEmail"))
      .matches(EMAIL_REGEX, {
        message: translate("invalidEmail"),
        excludeEmptyString: true,
      }),
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
    password: string()
      .required(translate("requiredPassword"))
      .min(
        PASSWORD_MIN_LENGTH,
        translate("minLengthPassword", { min: PASSWORD_MIN_LENGTH }),
      )
      .matches(PASSWORD_REGEX, translate("passwordComplexity")),
    confirmPassword: string()
      .required(translate("requiredConfirmPassword"))
      .oneOf([ref("password")], translate("passwordsDoNotMatch")),
  });
