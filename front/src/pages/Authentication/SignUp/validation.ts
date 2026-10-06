import type { TFunction } from "i18next";
import { boolean, object, ref, string } from "yup";
import { EMAIL_REGEX, PASSWORD_MIN_LENGTH, PASSWORD_REGEX } from "./constants";

export const getSchema = (translate: TFunction) =>
  object({
    email: string()
      .required(translate("requiredEmail"))
      .matches(EMAIL_REGEX, translate("invalidEmail")),
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
    termsAccepted: boolean()
      .oneOf([true], translate("termsAcceptedError"))
      .required(translate("termsAcceptedError")),
  });
