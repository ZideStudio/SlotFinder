import { FormProvider, useForm } from "react-hook-form";
import { useSignIn } from "./useSignIn";
import { Link } from "react-router";
import { yupResolver } from "@hookform/resolvers/yup";
import { useTranslation } from "react-i18next";
import { getSchema, type SignInFormType } from "./validation";
import { Button } from "@Front/ui/molecules/Button/Button";
import { TextField } from "@Front/components/fields/TextField/TextField";
import { InputErrorMessage } from "@Front/ui/atoms/Inputs/InputErrorMessage/InputErrorMessage";

import "./SignIn.scss";

export const SignIn = () => {
  const { t } = useTranslation("signIn");
  const { signIn, error } = useSignIn();
  const methods = useForm<SignInFormType>({
    resolver: yupResolver(getSchema(t)),
  });

  return (
    <FormProvider {...methods}>
      <form
        className="sign-in-form"
        onSubmit={methods.handleSubmit(signIn)}
        noValidate
      >
        <TextField
          name="identifier"
          label={t("identifierLabel")}
          autoComplete="username"
          required
        />
        <TextField
          name="password"
          label={t("passwordLabel")}
          type="password"
          autoComplete="current-password"
          required
        />
        {error && <InputErrorMessage>{t('error')}</InputErrorMessage>}

        <div className="sign-in-form__actions">
          <Button type="submit">{t("submitButton")}</Button>
          <Link className="sign-in-form__forgot" to="/forgot-password">
            {t("forgotPassword")}
          </Link>
        </div>
      </form>
    </FormProvider>
  );
};
