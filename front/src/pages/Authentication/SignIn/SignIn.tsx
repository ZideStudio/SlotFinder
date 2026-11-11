import { FormProvider, useForm } from "react-hook-form";
import { Link } from "react-router";
import { yupResolver } from "@hookform/resolvers/yup";
import { useTranslation } from "react-i18next";
import { getSchema, type SignInFormType } from "./validation";
import { Button } from "@Front/ui/molecules/Button/Button";
import { TextField } from "@Front/components/fields/TextField/TextField";

import "./SignIn.scss";

export const SignIn = () => {
  const { t } = useTranslation("signIn");
  const methods = useForm<SignInFormType>({
    resolver: yupResolver(getSchema(t)),
  });

  const onSubmit = (data: SignInFormType) => {
    // TODO: appel API de connexion
    console.log(data);
  };

  return (
    <FormProvider {...methods}>
      <form
        className="sign-in-form"
        onSubmit={methods.handleSubmit(onSubmit)}
        aria-labelledby="card-page-title"
        noValidate
      >
        <TextField
          name="identifier"
          label={t("identifierLabel")}
          autoComplete="username"
        />

        <TextField
          name="password"
          label={t("passwordLabel")}
          type="password"
          autoComplete="current-password"
        />

        <div className="sign-in__actions">
          <Button type="submit">{t("submitButton")}</Button>
          <Link className="sign-in__forgot" to="/forgot-password">
            {t("forgotPassword")}
          </Link>
        </div>
      </form>
    </FormProvider>
  );
};
