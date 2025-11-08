import { useMemo, useState, type ReactNode } from 'react';
import { AuthenticationContext } from './AuthenticationContext';
import { useCheckAuthentication } from './useCheckAuthentication';
import { usePostAuthentication } from './usePostAuthentication';

type AuthenticationContextProviderProps = {
  children: ReactNode;
};

export const AuthenticationContextProvider = ({ children }: AuthenticationContextProviderProps) => {
  const [isAuthenticated, setIsAuthenticated] = useState<boolean | undefined>();
  const { checkAuthentication } = useCheckAuthentication({
    onSuccess: () => setIsAuthenticated(true),
    onError: () => setIsAuthenticated(false),
  });
  const { postAuthRedirectPath, setPostAuthRedirectPath, resetPostAuthRedirectPath } = usePostAuthentication();

  const value = useMemo(
    () => ({
      isAuthenticated,
      postAuthRedirectPath,
      setPostAuthRedirectPath,
      resetPostAuthRedirectPath,
      checkAuthentication,
    }),
    [isAuthenticated, postAuthRedirectPath, setPostAuthRedirectPath, resetPostAuthRedirectPath, checkAuthentication],
  );

  return <AuthenticationContext.Provider value={value}>{children}</AuthenticationContext.Provider>;
};
