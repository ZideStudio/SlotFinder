import { postV1AccountForgotPassword } from "@Front/api/generated/account/account";
import type { AccountForgotPasswordDto } from "@Front/api/generated/slotFinderAPI.schemas";
import type { ErrorResponse } from "@Front/types/ErrorResponse";
import type { ErrorResponseCodeType } from "@Front/types/api.types";
import { useMutation } from "@tanstack/react-query";
import { useMemo } from "react";

type ForgotPasswordErrorCodeType = ErrorResponseCodeType<
  "PASSWORD_RESET_TOO_FREQUENT"
>;

type UseForgotPasswordReturn = {
  forgotPassword: (data: AccountForgotPasswordDto) => void;
  isLoading: boolean;
  isSuccess: boolean;
  errorCode?: ForgotPasswordErrorCodeType;
};

export const useForgotPassword = (): UseForgotPasswordReturn => {
  const mutation = useMutation<
    void,
    ErrorResponse<ForgotPasswordErrorCodeType>,
    AccountForgotPasswordDto
  >({
    mutationKey: ["forgotPassword"],
    mutationFn: (data) => postV1AccountForgotPassword(data),
  });

  const errorCode = useMemo(
    () => mutation.error?.getErrorCode(),
    [mutation.error],
  );

  return {
    forgotPassword: mutation.mutate,
    isLoading: mutation.isPending,
    isSuccess: mutation.isSuccess,
    errorCode,
  };
};