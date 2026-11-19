import type { TFunction } from "i18next";
import { type ObjectSchema } from "yup";
import type { WhoAreYouFormData } from "./types";
import { getCommonSchema } from "@Front/utils/helpers/schemas";

export const getSchema = (
  translate: TFunction<"whoAreYou">,
): ObjectSchema<WhoAreYouFormData> =>
  getCommonSchema(translate)
    .pick(["avatar", "username", "color", "termsAccepted"])
    .required();
