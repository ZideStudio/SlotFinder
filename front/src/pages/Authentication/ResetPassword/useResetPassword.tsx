import { postV1AccountResetPassword } from "@Front/api/generated/account/account";
import type { AccountResetPasswordDto } from "@Front/api/generated/slotFinderAPI.schemas";
import type { ErrorResponse } from "@Front/types/ErrorResponse";
import type { ErrorResponseCodeType } from "@Front/types/api.types";
import { useMutation } from "@tanstack/react-query";
import { useMemo } from "react";

type ResetPasswordErrorCodeType =
  ErrorResponseCodeType<"PASSWORD_RESET_TOO_FREQUENT">;

type UseResetPasswordReturn = {
  resetPassword: (data: AccountResetPasswordDto) => void;
  isLoading: boolean;
  isSuccess: boolean;
  errorCode?: ResetPasswordErrorCodeType;
};

export const useResetPassword = (): UseResetPasswordReturn => {
  const mutation = useMutation<
    void,
    ErrorResponse<ResetPasswordErrorCodeType>,
    AccountResetPasswordDto
  >({
    mutationKey: ["resetPassword"],
    mutationFn: (data) => postV1AccountResetPassword(data),
  });

  const errorCode = useMemo(
    () => mutation.error?.getErrorCode(),
    [mutation.error],
  );

  return {
    resetPassword: mutation.mutate,
    isLoading: mutation.isPending,
    isSuccess: mutation.isSuccess,
    errorCode,
  };
};
