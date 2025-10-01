import { authStatusApi } from '@Front/api/authentication/authStatusApi';
import type { AuthStatusResponseType } from '@Front/types/Authentication/authStatus/authStatus.types';
import type { AuthStatusErrorResponse } from '@Front/types/Authentication/authStatus/AuthStatusErrorResponse';
import { useMutation, type UseMutateAsyncFunction } from '@tanstack/react-query';
import { useLayoutEffect } from 'react';

type UseCheckAuthenticationProps = {
  onSuccess: () => void;
  onError: () => void;
};

export type UseCheckAuthenticationReturn = {
  checkAuthentication: UseMutateAsyncFunction<null, AuthStatusErrorResponse, void, unknown>;
};

export const useCheckAuthentication = ({
  onSuccess,
  onError,
}: UseCheckAuthenticationProps): UseCheckAuthenticationReturn => {
  const mutation = useMutation<AuthStatusResponseType, AuthStatusErrorResponse>({
    mutationFn: authStatusApi,
    retry: false,
    gcTime: 0,
    onSuccess,
    onError,
  });

  useLayoutEffect(() => {
    mutation.mutateAsync();
    // oxlint-disable-next-line exhaustive-deps
  }, []);

  return { checkAuthentication: mutation.mutateAsync };
};
