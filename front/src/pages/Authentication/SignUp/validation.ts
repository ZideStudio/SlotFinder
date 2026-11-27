import type { TFunction } from "i18next";
import { z } from "zod";
import { EMAIL_REGEX, PASSWORD_MIN_LENGTH, PASSWORD_REGEX } from "./constants";

export const getSchema = (translate: TFunction<"signUp">) =>
  z
    .object({
      email: z
        .string()
        .min(1, translate("requiredEmail"))
        .regex(EMAIL_REGEX, translate("invalidEmail")),
      password: z
        .string()
        .min(1, translate("requiredPassword"))
        .min(
          PASSWORD_MIN_LENGTH,
          translate("minLengthPassword", { min: PASSWORD_MIN_LENGTH }),
        )
        .regex(PASSWORD_REGEX, translate("passwordComplexity")),
      confirmPassword: z.string().min(1, translate("requiredConfirmPassword")),
      termsAccepted: z.custom<boolean>((value) => value === true, {
        message: translate("termsAcceptedError"),
      }),
    })
    .refine((form) => form.password === form.confirmPassword, {
      path: ["confirmPassword"],
      message: translate("passwordsDoNotMatch"),
    });
