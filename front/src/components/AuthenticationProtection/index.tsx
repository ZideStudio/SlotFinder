import { useAuthenticationContext } from '@Front/hooks/useAuthenticationContext';
import { appRoutes } from '@Front/routing/appRoutes';
import { useMemo, type ReactNode } from 'react';
import { Navigate, useLocation, useMatches, type UIMatch } from 'react-router';

type AuthenticationProtectionProps = {
  children: ReactNode;
};

export const AuthenticationProtection = ({ children }: AuthenticationProtectionProps) => {
  const { isAuthenticated } = useAuthenticationContext();
  const { pathname } = useLocation();
  const matches = useMatches() as UIMatch<unknown, { mustBeAuthenticate?: boolean }>[];

  const mustBeAuthenticate = useMemo(() => {
    const currentMatch = matches.find(match => match.pathname === pathname);
    return currentMatch?.handle?.mustBeAuthenticate;
  }, [matches, pathname]);

  if (isAuthenticated === undefined) {
    return null;
  }

  if (mustBeAuthenticate && !isAuthenticated) {
    return <Navigate to={appRoutes.signUp()} replace />;
  }

  if (mustBeAuthenticate === false && isAuthenticated) {
    return <Navigate to={appRoutes.dashboard()} replace />;
  }

  return children;
};
