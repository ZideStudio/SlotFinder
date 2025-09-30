import { signUpApi } from '@Front/api/signUpApi';
import type {
  SignUpErrorCodeType,
  SignUpFormType,
  SignUpRequestBodyType,
  SignUpResponseType,
} from '@Front/types/Authentication/signUp/signUp.types';
import type { SignUpErrorResponse } from '@Front/types/Authentication/signUp/SignUpErrorResponse';
import { useMutation } from '@tanstack/react-query';
import { useMemo } from 'react';

type UseSignUpApiReturn = {
  signUp: (userData: SignUpRequestBodyType) => Promise<SignUpResponseType>;
  isLoading: boolean;
  errorCode?: SignUpErrorCodeType;
};

export const useSignUp = (): UseSignUpApiReturn => {
  const mutation = useMutation<SignUpResponseType, SignUpErrorResponse, SignUpFormType>({
    mutationFn: signUpApi,
  });

  const errorCode = useMemo(() => mutation.error?.getErrorCode(), [mutation.error]);

  return {
    signUp: mutation.mutateAsync,
    isLoading: mutation.isPending,
    errorCode,
  };
};
