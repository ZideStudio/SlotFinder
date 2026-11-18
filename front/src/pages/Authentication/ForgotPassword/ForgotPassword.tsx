import { Button } from "@Front/ui/molecules/Button/Button";
import { yupResolver } from "@hookform/resolvers/yup";
import { FormProvider, useForm } from "react-hook-form";
import { useEffect } from "react";
import { useTranslation } from "react-i18next";
import "./ForgotPassword.scss";
import { getSchema } from "./validation";
import { TextField } from "@Front/components/fields/TextField/TextField";
import { useForgotPassword } from "./useForgotPassword";
import { useToastService } from "@Front/ui/utils/toast/hooks/useToastService";
import { InputErrorMessage } from "@Front/ui/atoms/Inputs/InputErrorMessage/InputErrorMessage";

type ForgotPasswordFormType = {
  email: string;
};

export const ForgotPassword = () => {
  const { forgotPassword, isLoading, isSuccess, errorCode } =
    useForgotPassword();
  const { t } = useTranslation("forgotPassword");
  const toastService = useToastService();
  const methods = useForm<ForgotPasswordFormType>({
    resolver: yupResolver(getSchema(t)),
  });

  useEffect(() => {
    if (isSuccess) {
      toastService.addToast(t("confirmation"), 5000);
    }
  }, [isSuccess, t, toastService]);

  return (
    <section>
      <FormProvider {...methods}>
        <form
          noValidate
          onSubmit={methods.handleSubmit(forgotPassword)}
          className="forgot-password-form"
        >
          <TextField name="email" label={t("email")} required />
          <Button type="submit" isLoading={isLoading}>
            {t("submit")}
          </Button>
        </form>
      </FormProvider>

      {errorCode && (
        <InputErrorMessage>{t(`error.${errorCode}`)}</InputErrorMessage>
      )}
    </section>
  );
};
