import { signUpApi } from '@Front/api/authentication/signUpApi';
import { useAuthenticationContext } from '@Front/hooks/useAuthenticationContext';
import type {
  SignUpErrorCodeType,
  SignUpFormType,
  SignUpResponseType,
} from '@Front/types/Authentication/signUp/signUp.types';
import type { SignUpErrorResponse } from '@Front/types/Authentication/signUp/SignUpErrorResponse';
import { useMutation } from '@tanstack/react-query';
import { useMemo } from 'react';
import { useTranslation } from 'react-i18next';

type UseSignUpApiReturn = {
  signUp: (userData: SignUpFormType) => void;
  isLoading: boolean;
  errorCode?: SignUpErrorCodeType;
};

export const useSignUp = (): UseSignUpApiReturn => {
  const { checkAuthentication } = useAuthenticationContext();
  const { i18n } = useTranslation();

  const mutation = useMutation<SignUpResponseType, SignUpErrorResponse, SignUpFormType>({
    mutationFn: ({ username, email, password }: SignUpFormType) => {
      return signUpApi({ username, email, password, language: i18n.language });
    },
    onSuccess: () => {
      checkAuthentication();
    },
  });

  const errorCode = useMemo(() => mutation.error?.getErrorCode(), [mutation.error]);

  return {
    signUp: mutation.mutate,
    isLoading: mutation.isPending,
    errorCode,
  };
};
