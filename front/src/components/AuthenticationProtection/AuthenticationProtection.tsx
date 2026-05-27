import { useAuthenticationContext } from "@Front/hooks/useAuthenticationContext";
import { appRoutes } from "@Front/routing/appRoutes";
import type { RouteHandle } from "@Front/routing/routeHandle";
import { useMemo, type ReactNode } from "react";
import { Navigate, useLocation, useMatches, type UIMatch } from "react-router";

type AuthenticationProtectionProps = {
  children: ReactNode;
};

export const AuthenticationProtection = ({
  children,
}: AuthenticationProtectionProps) => {
  const {
    isAuthenticated,
    postAuthRedirectPath,
    setPostAuthRedirectPath,
    resetPostAuthRedirectPath,
  } = useAuthenticationContext();
  const { pathname } = useLocation();
  const matches = useMatches() as UIMatch<unknown, RouteHandle>[];

  const mustBeAuthenticate = useMemo(() => {
    const currentMatch = matches.find((match) => match.pathname === pathname);

    if (currentMatch?.handle?.mustBeAuthenticate === true && !isAuthenticated) {
      setPostAuthRedirectPath(pathname);
    }

    if (
      currentMatch?.handle?.mustBeAuthenticate === false &&
      isAuthenticated &&
      postAuthRedirectPath
    ) {
      resetPostAuthRedirectPath();
    }

    return currentMatch?.handle?.mustBeAuthenticate;
  }, [
    pathname,
    matches,
    isAuthenticated,
    postAuthRedirectPath,
    setPostAuthRedirectPath,
    resetPostAuthRedirectPath,
  ]);

  if (isAuthenticated === undefined) {
    return null;
  }

  if (mustBeAuthenticate && !isAuthenticated) {
    return <Navigate to={appRoutes.signUp()} replace />;
  }

  if (mustBeAuthenticate === false && isAuthenticated) {
    return <Navigate to={postAuthRedirectPath ?? appRoutes.home()} replace />;
  }

  return children;
};
