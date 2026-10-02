import { useSignUp } from "@Front/pages/Authentication/SignUp/useSignUp";
import type { SignUpFormType } from "@Front/types/Authentication/signUp/signUp.types";
import { yupResolver } from "@hookform/resolvers/yup";
import { FormProvider, useForm } from "react-hook-form";
import { useTranslation } from "react-i18next";
import { Button } from "@Front/ui/molecules/Button/Button";
import { CheckboxField } from "@Front/components/fields/CheckboxField/CheckboxField";
import "./SignUp.scss";
import { OAuth } from "../OAuth/OAuth";
import { TextField } from "@Front/components/fields/TextField/TextField";
import { InputErrorMessage } from "@Front/ui/atoms/Inputs/InputErrorMessage/InputErrorMessage";
import { getSchema } from "./validation";

export const SignUp = () => {
  const { signUp, errorCode } = useSignUp();
  const { t } = useTranslation("signUp");
  const methods = useForm<SignUpFormType>({
    resolver: yupResolver(getSchema(t)),
  });

  return (
    <FormProvider {...methods}>
      <form onSubmit={methods.handleSubmit(signUp)}>
        <TextField
          name="email"
          label={t("email")}
          className="sign-up-input"
          id="email"
        />
        <TextField
          name="password"
          label={t("password")}
          className="sign-up-input"
          id="password"
          type="password"
        />
        <TextField
          name="confirmPassword"
          label={t("confirmPassword")}
          className="sign-up-input"
          id="confirmPassword"
          type="password"
        />
        <CheckboxField
          label={t("termsAcceptedLabel")}
          name="termsAccepted"
          required
        />
        <InputErrorMessage>
          {errorCode ? t(`error.${errorCode}`) : undefined}
        </InputErrorMessage>
        <hr />
        <OAuth />
        <Button type="submit">{t("submit")}</Button>
      </form>
    </FormProvider>
  );
};
