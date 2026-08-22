import { CardPage } from "@Front/components/CardPage/CardPage";
import { CheckboxField } from "@Front/components/fields/CheckboxField/CheckboxField";
import { ColorField } from "@Front/components/fields/ColorField/ColorField";
import { PictureUploadField } from "@Front/components/fields/PictureUploadField/PictureUploadField";
import { TextField } from "@Front/components/fields/TextField/TextField";
import { Button } from "@Front/ui/molecules/Button/Button";
import { yupResolver } from "@hookform/resolvers/yup";
import { FormProvider, useForm } from "react-hook-form";
import { useTranslation } from "react-i18next";
import { AVATAR_FILE_TYPES } from "./constants";
import type { WhoAreYouFormData } from "./types";
import { useWhoAreYou } from "./useWhoAreYou";
import { getSchema } from "./validation";

import { InputErrorMessage } from "@Front/ui/atoms/Inputs/InputErrorMessage/InputErrorMessage";
import "./WhoAreYou.scss";

export const WhoAreYou = () => {
  const { t } = useTranslation("whoAreYou");
  const methods = useForm<WhoAreYouFormData>({
    resolver: yupResolver(getSchema(t)),
  });
  const { handleSubmit, isLoading, submitError } = useWhoAreYou({
    setError: methods.setError,
  });

  return (
    <CardPage className="who-are-you" title={t("title")}>
      <FormProvider {...methods}>
        <form
          className="who-are-you-form"
          onSubmit={methods.handleSubmit(handleSubmit)}
          noValidate
        >
          <PictureUploadField
            label={t("avatarLabel")}
            name="avatar"
            accept={AVATAR_FILE_TYPES.join(",")}
            previewVariant="rounded"
            required
          />
          <TextField
            label={t("usernameLabel")}
            name="username"
            autoComplete="username"
            required
          />
          <ColorField
            label={t("colorLabel")}
            name="color"
            description={t("colorDescription")}
            required
          />
          <CheckboxField
            label={t("termsAcceptedLabel")}
            name="termsAccepted"
            required
          />
          <Button
            className="who-are-you-form__submit-button"
            type="submit"
            isLoading={isLoading}
          >
            Continuer
          </Button>
          <InputErrorMessage>{submitError}</InputErrorMessage>
        </form>
      </FormProvider>
    </CardPage>
  );
};
