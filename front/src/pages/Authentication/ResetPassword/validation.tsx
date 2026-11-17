import type { TFunction } from "i18next";
import { object, string } from "yup";

export const getSchema = (translate: TFunction) =>
  object({
    password: string()
      .required(translate("requiredPassword"))
      .min(8, translate("invalidPassword"))
      .max(100, translate("invalidPassword")),
    confirmPassword: string()
      .required(translate("requiredConfirmPassword"))
      .oneOf([string().nullable(), null], translate("passwordsDoNotMatch")),
  });
