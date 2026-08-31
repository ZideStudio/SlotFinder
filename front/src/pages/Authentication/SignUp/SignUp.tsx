import { useSignUp } from "@Front/pages/Authentication/SignUp/useSignUp";
import type { SignUpFormType } from "@Front/types/Authentication/signUp/signUp.types";
import { yupResolver } from "@hookform/resolvers/yup";
import { FormProvider, useForm } from "react-hook-form";
import { useTranslation } from "react-i18next";
import { getSchema } from "./validation";
import { TextInput } from "@Front/ui/molecules/Inputs/TextInput/TextInput";
import { Button } from "@Front/ui/molecules/Button/Button";
import { CardPage } from "@Front/components/CardPage/CardPage";
import { CheckboxField } from "@Front/components/fields/CheckboxField/CheckboxField";
import "./SignUp.scss";
import { OAuth } from "../OAuth/OAuth";

export const SignUp = () => {
  const { signUp } = useSignUp();
  const { t } = useTranslation("signUp");
  const methods = useForm<SignUpFormType>({
    resolver: yupResolver(getSchema(t)),
  });

  return (
    <CardPage title={t("title")} className="sign-up-page">
      <FormProvider {...methods}>
        <form onSubmit={methods.handleSubmit(signUp)}>
          <TextInput
            name="email"
            label={t("email")}
            id="email"
            error={methods.formState.errors.email?.message}
          />
          <TextInput
            name="password"
            label={t("password")}
            id="password"
            type="password"
            error={methods.formState.errors.password?.message}
          />
          <TextInput
            name="confirmPassword"
            label={t("confirmPassword")}
            id="confirmPassword"
            type="password"
            error={methods.formState.errors.confirmPassword?.message}
          />
          <CheckboxField
            label={t("termsAndConditions")}
            name="termsAndConditions"
            id="termsAndConditions"
            required
          />
          <hr />
          <OAuth />
          <Button type="submit">{t("submit")}</Button>
        </form>
      </FormProvider>
    </CardPage>
    // <FormProvider {...methods}>
    //   <form
    //     OnSubmit={methods.handleSubmit(signUp)}
    //     Style={{
    //       MaxWidth: 400,
    //       Margin: "0 auto",
    //       Display: "flex",
    //       FlexDirection: "column",
    //       Gap: "1.5rem",
    //     }}
    //     Aria-labelledby="signup-legend"
    //   >
    //     <fieldset style={{ border: "none", padding: 0, margin: 0 }}>
    //       <legend
    //         Id="signup-legend"
    //         Style={{ fontWeight: "bold", marginBottom: "1rem" }}
    //       >
    //         {t("title")}
    //       </legend>

    //       <div
    //         Style={{ display: "flex", flexDirection: "column", gap: "0.5rem" }}
    //       >
    //         <TextInput
    //           Label={t("username")}
    //           Id="username"
    //           Type="text"
    //           AutoComplete="username"
    //           Aria-describedby={
    //             Methods.formState.errors.username ? "username-error" : undefined
    //           }
    //           {...methods.register("username", { required: true })}
    //         />
    //         {methods.formState.errors.username && (
    //           <span
    //             Id="username-error"
    //             Role="alert"
    //             Style={{ color: "red", marginTop: 2 }}
    //           >
    //             {methods.formState.errors.username.message}
    //           </span>
    //         )}
    //       </div>

    //       <div
    //         Style={{ display: "flex", flexDirection: "column", gap: "0.5rem" }}
    //       >
    //         <TextInput
    //           Label={t("email")}
    //           Id="email"
    //           Type="email"
    //           AutoComplete="email"
    //           Aria-describedby={
    //             Methods.formState.errors.email ? "email-error" : undefined
    //           }
    //           {...methods.register("email", { required: true })}
    //         />
    //         {methods.formState.errors.email && (
    //           <span
    //             Id="email-error"
    //             Role="alert"
    //             Style={{ color: "red", marginTop: 2 }}
    //           >
    //             {methods.formState.errors.email.message}
    //           </span>
    //         )}
    //       </div>

    //       <div
    //         Style={{ display: "flex", flexDirection: "column", gap: "0.5rem" }}
    //       >
    //         <TextInput
    //           Label={t("password")}
    //           Id="password"
    //           Type="password"
    //           AutoComplete="new-password"
    //           Aria-describedby={
    //             Methods.formState.errors.password ? "password-error" : undefined
    //           }
    //           {...methods.register("password", { required: true })}
    //         />
    //         {methods.formState.errors.password && (
    //           <span
    //             Id="password-error"
    //             Role="alert"
    //             Style={{ color: "red", marginTop: 2 }}
    //           >
    //             {methods.formState.errors.password.message}
    //           </span>
    //         )}
    //       </div>

    //       <div
    //         Style={{ display: "flex", flexDirection: "column", gap: "0.5rem" }}
    //       >
    //         <TextInput
    //           Label={t("confirmPassword")}
    //           Id="confirmPassword"
    //           Type="password"
    //           AutoComplete="new-password"
    //           Aria-describedby={
    //             Methods.formState.errors.confirmPassword
    //               ? "confirmPassword-error"
    //               : undefined
    //           }
    //           {...methods.register("confirmPassword", { required: true })}
    //         />
    //         {methods.formState.errors.confirmPassword && (
    //           <span
    //             Id="confirmPassword-error"
    //             Role="alert"
    //             Style={{ color: "red", marginTop: 2 }}
    //           >
    //             {methods.formState.errors.confirmPassword.message}
    //           </span>
    //         )}
    //       </div>
    //     </fieldset>
    //     {errorCode && (
    //       <span role="alert" style={{ color: "red" }}>
    //         {t(`error.${errorCode}`)}
    //       </span>
    //     )}
    //     <Button type="submit" disabled={isLoading}>
    //       {t("submit")}
    //     </Button>
    //   </form>
    // </FormProvider>
  );
};
