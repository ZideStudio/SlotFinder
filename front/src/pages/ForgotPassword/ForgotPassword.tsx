import { postV1AccountForgotPassword } from "@Front/api/generated/account/account";
import { signUpRoutes } from "@Front/pages/Authentication/SignUp";
import type { ErrorResponse } from "@Front/types/ErrorResponse";
import type { ErrorResponseCodeType } from "@Front/types/api.types";
import { Button } from "@Front/ui/molecules/Button/Button";
import { TextInput } from "@Front/ui/molecules/Inputs/TextInput/TextInput";
import { yupResolver } from "@hookform/resolvers/yup";
import { useMutation } from "@tanstack/react-query";
import { useForm } from "react-hook-form";
import { useTranslation } from "react-i18next";
import { Link } from "react-router";
import { object, string } from "yup";
import "./ForgotPassword.scss";

type ForgotPasswordFormType = {
  email: string;
};

type ForgotPasswordErrorCodeType =
  ErrorResponseCodeType<"PASSWORD_RESET_TOO_FREQUENT">;

export const ForgotPassword = () => {
  const { t } = useTranslation("forgotPassword");
  const schema = object({
    email: string().required(t("requiredEmail")).email(t("invalidEmail")),
  });
  const mutation = useMutation<
    void,
    ErrorResponse<ForgotPasswordErrorCodeType>,
    ForgotPasswordFormType
  >({
    mutationKey: ["forgotPassword"],
    mutationFn: ({ email }) => postV1AccountForgotPassword({ email }),
  });
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<ForgotPasswordFormType>({
    resolver: yupResolver(schema),
  });

  return (
    <section className="forgot-password">
      <h2 className="forgot-password__title">{t("title")}</h2>
      {mutation.isSuccess ? (
        <div className="forgot-password__confirmation">
          <output>{t("confirmation")}</output>
          <Button
            type="button"
            variant="secondary"
            onClick={() => mutation.reset()}
          >
            {t("tryAgain")}
          </Button>
        </div>
      ) : (
        <>
          <p className="forgot-password__description">{t("description")}</p>
          <form
            className="forgot-password__form"
            onSubmit={handleSubmit((values) => mutation.mutate(values))}
            noValidate
          >
            <TextInput
              {...register("email")}
              id="forgot-password-email"
              label={t("email")}
              type="email"
              autoComplete="email"
              required
              error={errors.email?.message}
            />
            {mutation.error && (
              <p className="forgot-password__error" role="alert">
                {t(`error.${mutation.error.getErrorCode()}`)}
              </p>
            )}
            <Button type="submit" isLoading={mutation.isPending}>
              {t("submit")}
            </Button>
          </form>
        </>
      )}
      <p className="forgot-password__signup">
        {t("noAccount")}{" "}
        <Link to={`/${signUpRoutes.path}`}>{t("createAccount")}</Link>
      </p>
    </section>
  );
};
