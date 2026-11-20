import type { TFunction } from "i18next";
import { type ObjectSchema } from "yup";
import type { WhoAreYouFormData } from "./types";
import { getUserSchema } from "@Front/utils/validation/userSchema";

export const getSchema = (
  translate: TFunction<"whoAreYou">,
): ObjectSchema<WhoAreYouFormData> =>
  getUserSchema(translate)
    .pick(["avatar", "username", "color", "termsAccepted"])
    .required();
