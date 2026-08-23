import { patchV1Account } from "@Front/api/generated/account/account";
import type {
  AccountAccountResponseDto,
  AccountAccountUpdateDto,
} from "@Front/api/generated/slotFinderAPI.schemas";
import { type ErrorResponse } from "@Front/types/ErrorResponse";
import { useMutation } from "@tanstack/react-query";
import { useMemo } from "react";
import type { PatchAccountErrorCodeType } from "./types";

type UsePatchAccountProps = {
  onError?: (error: ErrorResponse<PatchAccountErrorCodeType>) => void;
};

type UsePatchAccountReturn = {
  patchAccount: (data: AccountAccountUpdateDto) => void;
  isLoading: boolean;
  errorCode?: PatchAccountErrorCodeType;
};

export const usePatchAccount = ({
  onError,
}: UsePatchAccountProps = {}): UsePatchAccountReturn => {
  const mutation = useMutation<
    AccountAccountResponseDto,
    ErrorResponse<PatchAccountErrorCodeType>,
    AccountAccountUpdateDto
  >({
    mutationKey: ["patchAccount"],
    mutationFn: (data) => patchV1Account(data),
    onError,
  });

  const errorCode = useMemo(
    () => mutation.error?.getErrorCode(),
    [mutation.error],
  );

  return {
    patchAccount: mutation.mutateAsync,
    isLoading: mutation.isPending,
    errorCode,
  };
};
