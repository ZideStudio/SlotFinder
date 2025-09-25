import { signUpApi } from '@Front/api/signUpApi';
import type { SignUpFormType, SignUpRequestBodyType, SignUpResponseType } from '@Front/types/Authentication/signUp.types';
import { getFormattedErrorMessage } from '@Front/utils/getFormattedErrorMessage';
import { useMutation } from '@tanstack/react-query';

type UseSignUpApiReturn = {
  signUp: (userData: SignUpRequestBodyType) => Promise<SignUpResponseType>;
  isLoading: boolean;
  error?: string;
};

export const useSignUp = (): UseSignUpApiReturn => {
  const mutation = useMutation<SignUpResponseType, Error, SignUpFormType>({
    mutationFn: signUpApi,
  });

  return {
    signUp: mutation.mutateAsync,
    isLoading: mutation.isPending,
    error: getFormattedErrorMessage(mutation.error),
  };
};
