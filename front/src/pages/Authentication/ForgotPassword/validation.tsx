import { getCommonSchema } from "@Front/utils/helpers/schemas";
import type { TFunction } from "i18next";
import { StringSchema } from "yup";

export const getSchema = (translate: TFunction) => {
  const commonSchema = getCommonSchema(translate);

  return commonSchema
    .shape({
      email: (commonSchema.fields.email as StringSchema).required(
        translate("requiredEmail"),
      ),
    })
    .pick(["email"])
    .required();
};
