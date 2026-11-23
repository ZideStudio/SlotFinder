import { signUpRoutes } from "@Front/pages/Authentication/SignUp";
import { Button } from "@Front/ui/molecules/Button/Button";
import { yupResolver } from "@hookform/resolvers/yup";
import { FormProvider, useForm } from "react-hook-form";
import { useEffect } from "react";
import { useTranslation } from "react-i18next";
import { Link } from "react-router";
import "./ForgotPassword.scss";
import { getSchema } from "./validation";
import { TextField } from "@Front/components/fields/TextField/TextField";
import { InputErrorMessage } from "@Front/ui/atoms/Inputs/InputErrorMessage/InputErrorMessage";
import { useToastService } from "@Front/ui/utils/toast/hooks/useToastService";
import { useResetPassword } from "./useResetPassword";

type ResetPasswordFormType = {
  password: string;
  confirmPassword: string;
  token: string;
};

export const ResetPassword = () => {
  const { resetPassword, isLoading, errorCode } =
    useResetPassword();
  const { t } = useTranslation("resetPassword");
  const methods = useForm<ResetPasswordFormType>({
    resolver: yupResolver(getSchema(t)),
  });

  return (
    <section className="reset-password">
      <FormProvider {...methods}>
        <form
          noValidate
          onSubmit={methods.handleSubmit(resetPassword)}
          className="reset-password__form"
        >
          <TextField name="password" label={t("password")} type="password" required />
          <TextField name="confirmPassword" label={t("confirmPassword")} type="password" required />
          <Button type="submit" isLoading={isLoading}>
            {t("submit")}
          </Button>
        </form>
        {errorCode && (
          <InputErrorMessage>{t(`error.${errorCode}`)}</InputErrorMessage>
        )}
        <p className="reset-password__signup">
          {t("noAccount")}{" "}
          <Link to={`/${signUpRoutes.path}`}>{t("createAccount")}</Link>
        </p>
      </FormProvider>
    </section>
  );
};
