import { type TFunction } from "i18next";
import { object, string, type InferType } from "yup";

export const getSchema = (translate: TFunction<"signIn">) =>
  object({
    identifier: string().required(translate("requiredIdentifier")),
    password: string().required(translate("requiredPassword")),
  });

export type SignInFormType = InferType<ReturnType<typeof getSchema>>;
