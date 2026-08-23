import { getV1AuthStatus } from "@Front/api/generated/authentication/authentication";
import type { AuthStatusErrorCodeType } from "@Front/types/Authentication/authStatus/authStatus.types";
import { type ErrorResponse } from "@Front/types/ErrorResponse";
import { useMutation, type UseMutateFunction } from "@tanstack/react-query";
import { useLayoutEffect } from "react";

type UseCheckAuthenticationProps = {
  onSuccess: () => void;
  onError: (
    error: ErrorResponse<AuthStatusErrorCodeType>,
  ) => Promise<unknown> | unknown;
};

export type UseCheckAuthenticationReturn = {
  checkAuthentication: UseMutateFunction<
    void,
    ErrorResponse<AuthStatusErrorCodeType>,
    void,
    unknown
  >;
};

export const useCheckAuthentication = ({
  onSuccess,
  onError,
}: UseCheckAuthenticationProps): UseCheckAuthenticationReturn => {
  const mutation = useMutation<void, ErrorResponse<AuthStatusErrorCodeType>>({
    mutationKey: ["checkAuthentication"],
    mutationFn: () => getV1AuthStatus(),
    retry: false,
    gcTime: 0,
    onSuccess,
    onError,
  });

  useLayoutEffect(() => {
    mutation.mutate();
    // oxlint-disable-next-line exhaustive-deps
  }, []);

  return { checkAuthentication: mutation.mutate };
};
