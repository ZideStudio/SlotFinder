import type { TFunction } from "i18next";
import { object, string } from "yup";

const EMAIL_REGEX = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/u;

export const getSchema = (translate: TFunction) =>
  object({
    email: string()
      .required(translate("requiredEmail"))
      .matches(EMAIL_REGEX, translate("invalidEmail"))
    });
