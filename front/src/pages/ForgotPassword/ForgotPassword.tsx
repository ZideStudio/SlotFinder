import { signUpRoutes } from "@Front/pages/Authentication/SignUp";
import { Button } from "@Front/ui/molecules/Button/Button";
import { yupResolver } from "@hookform/resolvers/yup";
import { FormProvider, useForm } from "react-hook-form";
import { useTranslation } from "react-i18next";
import { Link } from "react-router";
import "./ForgotPassword.scss";
import { getSchema } from "./validation";
import { TextField } from "@Front/components/fields/TextField/TextField";
import { InputErrorMessage } from "@Front/ui/atoms/Inputs/InputErrorMessage/InputErrorMessage";
import { useForgotPassword } from "./useForgotPassword";

type ForgotPasswordFormType = {
  email: string;
};

export const ForgotPassword = () => {
  const { forgotPassword, isLoading, errorCode } =
    useForgotPassword();
  const { t } = useTranslation("forgotPassword");
  const methods = useForm<ForgotPasswordFormType>({
    resolver: yupResolver(getSchema(t)),
  });

  return (
    <div className="forgot-password">
      <FormProvider {...methods}>
        <form
          noValidate
          onSubmit={methods.handleSubmit(forgotPassword)}
          className="sign-up-form"
        >
          <TextField
            name="email"
            label={t("email")}
            required
            className="sign-up-input"
          />
          <Button
            type="submit"
            isLoading={isLoading}
            className="sign-up-button"
          >
            {t("submit")}
          </Button>
        </form>
        {errorCode && (
          <InputErrorMessage>{t(`error.${errorCode}`)}</InputErrorMessage>
        )}
        <p className="forgot-password__signup">
          {t("noAccount")}{" "}
          <Link to={`/${signUpRoutes.path}`}>{t("createAccount")}</Link>
        </p>
      </FormProvider>
    </div>
  );
};
