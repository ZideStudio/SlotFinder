import { usePatchAccount } from "@Front/api/account/patchAccount/usePatchAccount";
import { usePatchAccountAvatar } from "@Front/api/account/patchAccountAvatar/usePatchAccountAvatar";
import type { AccountAccountUpdateDto } from "@Front/api/generated/slotFinderAPI.schemas";
import { useCallback, useState } from "react";
import type { UseFormSetError } from "react-hook-form";
import { useTranslation } from "react-i18next";
import type { WhoAreYouFormData } from "./types";

type UseWhoAreYouProps = {
  setError: UseFormSetError<WhoAreYouFormData>;
};

type UseWhoAreYouReturn = {
  handleSubmit: (data: WhoAreYouFormData) => Promise<void>;
  isLoading: boolean;
  submitError: string | null;
};

export const useWhoAreYou = ({
  setError,
}: UseWhoAreYouProps): UseWhoAreYouReturn => {
  const { t } = useTranslation("whoAreYou");
  const [submitError, setSubmitError] = useState<string | null>(null);

  const { patchAccount, isLoading: isAccountLoading } = usePatchAccount({
    onError: (error) => {
      const errorCode = error.getErrorCode();

      if (errorCode === "USERNAME_ALREADY_TAKEN") {
        setError("username", {
          message: t("error.USERNAME_ALREADY_TAKEN"),
        });
        return;
      }

      if (errorCode === "INVALID_COLOR_FORMAT") {
        setError("color", { message: t("error.INVALID_COLOR_FORMAT") });
        return;
      }

      setSubmitError(t("error.SERVER_ERROR"));
    },
  });

  const { patchAccountAvatar, isLoading: isAvatarLoading } =
    usePatchAccountAvatar({
      onError: () => {
        setError("avatar", { message: t("error.AVATAR_UPLOAD_FAILED") });
      },
    });

  const handleSubmit = useCallback(
    async (data: WhoAreYouFormData) => {
      setSubmitError(null);

      const patchAccountDTO: AccountAccountUpdateDto = {
        username: data.username,
        color: data.color,
        termsAccepted: data.termsAccepted,
        timeZone: Temporal.Now.timeZoneId(),
      };

      await Promise.all([
        patchAccount(patchAccountDTO),
        patchAccountAvatar({ image: data.avatar[0] }),
      ]);
    },
    [patchAccount, patchAccountAvatar],
  );

  return {
    handleSubmit,
    isLoading: isAccountLoading || isAvatarLoading,
    submitError,
  };
};
