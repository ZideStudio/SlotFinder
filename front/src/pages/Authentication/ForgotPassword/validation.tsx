import type { TFunction } from "i18next";
import { z } from "zod";

const EMAIL_REGEX = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/u;

export const getSchema = (translate: TFunction<"forgotPassword">) =>
  z.object({
    email: z
      .string()
      .min(1, translate("requiredEmail"))
      .max(320, translate("invalidEmail"))
      .regex(EMAIL_REGEX, translate("invalidEmail")),
  });
