import { useSignUp } from "@Front/pages/Authentication/SignUp/useSignUp";
import type { SignUpFormType } from "@Front/types/Authentication/signUp/signUp.types";
import { yupResolver } from "@hookform/resolvers/yup";
import { FormProvider, useForm } from "react-hook-form";
import { useTranslation } from "react-i18next";
import { Button } from "@Front/ui/molecules/Button/Button";
import { CheckboxField } from "@Front/components/fields/CheckboxField/CheckboxField";
import "./SignUp.scss";
import { TextField } from "@Front/components/fields/TextField/TextField";
import { InputErrorMessage } from "@Front/ui/atoms/Inputs/InputErrorMessage/InputErrorMessage";
import { getSchema } from "./validation";
import { Link } from "@Front/ui/atoms/Link/Link";

export const SignUp = () => {
  const { signUp, isLoading, errorCode } = useSignUp();
  const { t } = useTranslation("signUp");
  const methods = useForm<SignUpFormType>({
    resolver: yupResolver(getSchema(t)),
  });

  return (
    <FormProvider {...methods}>
      <form noValidate onSubmit={methods.handleSubmit(signUp)} className="sign-up-form">
        <TextField name="email" label={t("email")} className="sign-up-input" />
        <TextField
          name="password"
          label={t("password")}
          className="sign-up-input"
          type="password"
        />
        <TextField
          name="confirmPassword"
          label={t("confirmPassword")}
          className="sign-up-input"
          type="password"
        />
        <CheckboxField
          label={t("termsAcceptedLabel")}
          name="termsAccepted"
          required
        />
        {errorCode && (
          <InputErrorMessage>{t(`error.${errorCode}`)}</InputErrorMessage>
        )}
        <Button type="submit" isLoading={isLoading}>
          {t("submit")}
        </Button>
      </form>
      {/* <p>
        {t("existingAccount")} <Link as="a" href="/login">{t("login")}</Link>
      </p> */}
    </FormProvider>
    
  );
};
