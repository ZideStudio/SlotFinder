import type { TFunction } from "i18next";
import { getUserSchema } from "@Front/utils/validation/userSchema";

export const getSchema = (translate: TFunction<"signUp">) => {
  const commonSchema = getUserSchema(translate);

  return commonSchema.pick([
    "email",
    "password",
    "confirmPassword",
    "termsAccepted",
  ]);
};
