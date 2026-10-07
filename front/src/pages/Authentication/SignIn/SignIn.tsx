import { FormProvider, useForm } from "react-hook-form";
import { Link } from "react-router";
import { yupResolver } from "@hookform/resolvers/yup";
import { useTranslation } from "react-i18next";
import { getSchema, type SignInFormType } from "./validation";
import { Heading } from "@Front/ui/atoms/Heading/Heading";
import { TextInput } from "@Front/ui/molecules/Inputs/TextInput/TextInput";
import { Button } from "@Front/ui/molecules/Button/Button";
import { OAuth } from "../OAuth/OAuth";

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
    <div className="sign-in">
      <Heading id="signin-title" className="sign-in__title" level={1}>
        {t("title")}
      </Heading>

      <div className="sign-in__card">
        <FormProvider {...methods}>
          <form
            className="sign-in__form"
            onSubmit={methods.handleSubmit(onSubmit)}
            aria-labelledby="signin-title"
            noValidate
          >
            <TextInput
              id="identifier"
              label={t("identifierLabel")}
              autoComplete="username"
              {...methods.register("identifier")}
            />

            <TextInput
              id="password"
              label={t("passwordLabel")}
              type="password"
              autoComplete="current-password"
              {...methods.register("password")}
            />

            <div className="sign-in__actions">
              <Button type="submit">{t("submitButton")}</Button>
              <Link className="sign-in__forgot" to="/forgot-password">
                {t("forgotPassword")}
              </Link>
            </div>
          </form>
        </FormProvider>

        <hr className="sign-in__divider" />

        <OAuth />

        <p className="sign-in__signup">
          {t("noAccount")} <Link to="/sign-up">{t("signUp")}</Link>
        </p>
      </div>
    </div>
  );
};
