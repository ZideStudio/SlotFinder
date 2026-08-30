import type { AuthStatusErrorCodeType } from "@Front/types/Authentication/authStatus/authStatus.types";
import { useMemo, useState, type ReactNode } from "react";
import { AuthenticationContext } from "./AuthenticationContext";
import { useCheckAuthentication } from "./useCheckAuthentication";
import { usePostAuthentication } from "./usePostAuthentication";

type AuthenticationContextProviderProps = {
  children: ReactNode;
};

export const AuthenticationContextProvider = ({
  children,
}: AuthenticationContextProviderProps) => {
  const [isAuthenticated, setIsAuthenticated] = useState<boolean | undefined>();
  const [authenticationError, setAuthenticationError] =
    useState<AuthStatusErrorCodeType>();
  const { checkAuthentication } = useCheckAuthentication({
    onSuccess: () => {
      setIsAuthenticated(true);
      setAuthenticationError(undefined);
    },
    onError: (error) => {
      setIsAuthenticated(false);
      setAuthenticationError(error.getErrorCode());
    },
  });
  const {
    postAuthRedirectPath,
    setPostAuthRedirectPath,
    resetPostAuthRedirectPath,
  } = usePostAuthentication();

  const value = useMemo(
    () => ({
      isAuthenticated,
      authenticationError,
      postAuthRedirectPath,
      setPostAuthRedirectPath,
      resetPostAuthRedirectPath,
      checkAuthentication,
    }),
    [
      isAuthenticated,
      authenticationError,
      postAuthRedirectPath,
      setPostAuthRedirectPath,
      resetPostAuthRedirectPath,
      checkAuthentication,
    ],
  );

  return (
    <AuthenticationContext.Provider value={value}>
      {children}
    </AuthenticationContext.Provider>
  );
};
