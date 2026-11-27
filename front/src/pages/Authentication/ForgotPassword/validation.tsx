import { getUserSchema } from "@Front/utils/validation/userSchema";
import type { TFunction } from "i18next";

export const getSchema = (translate: TFunction<"forgotPassword">) => {
  const commonSchema = getUserSchema(translate);

  return commonSchema.pick(["email"]);
};
