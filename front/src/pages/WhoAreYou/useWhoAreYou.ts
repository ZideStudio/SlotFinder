import {
  getAccountMeQueryKey,
  useGetAccountMe,
} from "@Front/api/account/getAccountMe/useGetAccountMe";
import type { PatchAccountErrorCodeType } from "@Front/api/account/patchAccount/types";
import { usePatchAccount } from "@Front/api/account/patchAccount/usePatchAccount";
import { usePatchAccountAvatar } from "@Front/api/account/patchAccountAvatar/usePatchAccountAvatar";
import type { ErrorResponse } from "@Front/types/ErrorResponse";
import { urlToFileList } from "@Front/utils/helpers/urlToFileList";
import { useQueryClient } from "@tanstack/react-query";
import { useCallback, useEffect, useRef, useState } from "react";
import type { UseFormReset, UseFormSetError } from "react-hook-form";
import { useTranslation } from "react-i18next";
import type { WhoAreYouFormData } from "./types";
import { buildAccountUpdateDTO } from "./whoAreYou.helpers";

type UseWhoAreYouProps = {
  setError: UseFormSetError<WhoAreYouFormData>;
  reset: UseFormReset<WhoAreYouFormData>;
};

type UseWhoAreYouReturn = {
  handleSubmit: (data: WhoAreYouFormData) => void;
  isSubmitting: boolean;
  submitError: string | null;
  defaultAvatarUrl?: string;
};

export const useWhoAreYou = ({
  setError,
  reset,
}: UseWhoAreYouProps): UseWhoAreYouReturn => {
  const { t } = useTranslation("whoAreYou");
  const [submitError, setSubmitError] = useState<string | null>(null);
  const hasInitialized = useRef(false);
  const { data: accountData } = useGetAccountMe();
  const queryClient = useQueryClient();

  const handlePatchAccountError = useCallback(
    (error: ErrorResponse<PatchAccountErrorCodeType>) => {
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
    [setError, t],
  );

  const handlePatchAvatarError = useCallback(() => {
    setError("avatar", { message: t("error.AVATAR_UPLOAD_FAILED") });
  }, [setError, t]);

  const { patchAccount, isLoading: isAccountPatchLoading } = usePatchAccount({
    onError: handlePatchAccountError,
  });

  const { patchAccountAvatar, isLoading: isAvatarPatchLoading } =
    usePatchAccountAvatar({
      onError: handlePatchAvatarError,
    });

  // Initialize form values with account data when the component mounts
  useEffect(() => {
    if (hasInitialized.current || !accountData) {
      return;
    }

    const initializeFormValues = async () => {
      const avatarFileList = accountData.avatarUrl
        ? await urlToFileList(accountData.avatarUrl, "avatar.jpg")
        : undefined;

      reset({
        avatar: avatarFileList,
        username: accountData.username,
        color: accountData.color,
      });

      hasInitialized.current = true;
    };

    initializeFormValues();
  }, [accountData, reset]);

  const handleSubmit = useCallback(
    async (data: WhoAreYouFormData) => {
      setSubmitError(null);

      const patchAccountDTO = buildAccountUpdateDTO(data);
      const [avatarFile] = data.avatar;

      try {
        await Promise.all([
          patchAccount(patchAccountDTO),
          patchAccountAvatar({ image: avatarFile }),
        ]);

        queryClient.invalidateQueries({ queryKey: getAccountMeQueryKey });
      } catch {
        setSubmitError(t("error.SERVER_ERROR"));
      }
    },
    [patchAccount, patchAccountAvatar, queryClient, t],
  );

  return {
    handleSubmit,
    isSubmitting: isAccountPatchLoading || isAvatarPatchLoading,
    submitError,
    defaultAvatarUrl: accountData?.avatarUrl,
  };
};
