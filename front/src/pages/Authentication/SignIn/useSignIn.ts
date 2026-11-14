import { signInApi } from "@Front/api/authentication/signInApi";
import { useAuthenticationContext } from "@Front/hooks/useAuthenticationContext";
import { type ErrorResponse } from "@Front/types/ErrorResponse";
import { useMutation } from "@tanstack/react-query";
import { type SignInFormType } from "./validation";

type UseSignInApiReturn = {
  signIn: (userData: SignInFormType) => void;
  isLoading: boolean;
  error: ErrorResponse | null;
};

export const useSignIn = (): UseSignInApiReturn => {
  const { checkAuthentication } = useAuthenticationContext();

  const mutation = useMutation<void, ErrorResponse, SignInFormType>({
    mutationKey: ["signIn"],
    mutationFn: async ({ identifier, password }: SignInFormType) => {
      await signInApi({
        identifier,
        password,
      });
    },
    onSuccess: () => {
      checkAuthentication();
    },
  });

  return {
    signIn: mutation.mutate,
    isLoading: mutation.isPending,
    error: mutation.error ?? null,
  };
};
