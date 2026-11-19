import type { TFunction } from "i18next";
import type { StringSchema } from "yup";
import { getCommonSchema } from "@Front/utils/helpers/schemas";

export const getSchema = (translate: TFunction) => {
  const commonSchema = getCommonSchema(translate);

  return commonSchema
    .shape({
      email: (commonSchema.fields.email as StringSchema).required(
        translate("requiredEmail"),
      ),
    })
    .pick(["email", "password", "confirmPassword", "termsAccepted"])
    .required();
};
