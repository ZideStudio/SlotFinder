import { useMemo, useState, type ReactNode } from 'react';
import { AuthenticationContext } from './AuthenticationContext';
import { useCheckAuthentication } from './useCheckAuthentication';

type AuthenticationContextProviderProps = {
  children: ReactNode;
};

export const AuthenticationContextProvider = ({ children }: AuthenticationContextProviderProps) => {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const { checkAuthentication } = useCheckAuthentication({
    onSuccess: () => setIsAuthenticated(true),
    onError: () => setIsAuthenticated(false),
  });

  const value = useMemo(() => ({ isAuthenticated, checkAuthentication }), [isAuthenticated, checkAuthentication]);

  return <AuthenticationContext.Provider value={value}>{children}</AuthenticationContext.Provider>;
};
