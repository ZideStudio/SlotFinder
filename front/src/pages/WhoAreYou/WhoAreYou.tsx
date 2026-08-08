import { CardPage } from "@Front/components/CardPage/CardPage";
import { CheckboxField } from "@Front/components/fields/CheckboxField/CheckboxField";
import { ColorField } from "@Front/components/fields/ColorField/ColorField";
import { PictureUploadField } from "@Front/components/fields/PictureUploadField/PictureUploadField";
import { TextField } from "@Front/components/fields/TextField/TextField";
import { Button } from "@Front/ui/molecules/Button/Button";
import { FormProvider, useForm } from "react-hook-form";
import { useTranslation } from "react-i18next";

import "./WhoAreYou.scss";

export const WhoAreYou = () => {
  const { t } = useTranslation("whoAreYou");
  const methods = useForm();

  return (
    <CardPage className="who-are-you" title={t("title")}>
      <FormProvider {...methods}>
        <form
          className="who-are-you-form"
          onSubmit={methods.handleSubmit((data) => console.log(data))}
        >
          <PictureUploadField label={t("avatarLabel")} name="avatar" required />
          <TextField
            label={t("userNameLabel")}
            name="userName"
            required
            autoComplete="username"
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
          <Button className="who-are-you-form__submit-button" type="submit">
            Continuer
          </Button>
        </form>
      </FormProvider>
    </CardPage>
  );
};
