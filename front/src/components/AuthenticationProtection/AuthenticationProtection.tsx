import { useAuthenticationContext } from "@Front/hooks/useAuthenticationContext";
import LoaderPage from "@Front/pages/Loader/LoaderPage";
import { appRoutes } from "@Front/routing/appRoutes";
import { useEffect, type ReactNode } from "react";
import { Navigate, useLocation, useMatches, type UIMatch } from "react-router";
import type { RouteHandle } from "@Front/routing/routeHandle";

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

  const currentMatch = matches.at(-1);
  const mustBeAuthenticate = currentMatch?.handle?.mustBeAuthenticate;

  useEffect(() => {
    if (mustBeAuthenticate === true && !isAuthenticated) {
      setPostAuthRedirectPath(pathname);
      return;
    }

    if (
      mustBeAuthenticate === false &&
      isAuthenticated &&
      postAuthRedirectPath
    ) {
      resetPostAuthRedirectPath();
    }
  }, [
    mustBeAuthenticate,
    isAuthenticated,
    pathname,
    postAuthRedirectPath,
    setPostAuthRedirectPath,
    resetPostAuthRedirectPath,
  ]);

  if (isAuthenticated === undefined) {
    return <LoaderPage />;
  }

  if (mustBeAuthenticate && !isAuthenticated) {
    return <Navigate to={appRoutes.signUp()} replace />;
  }

  if (mustBeAuthenticate === false && isAuthenticated) {
    return <Navigate to={postAuthRedirectPath ?? appRoutes.home()} replace />;
  }

  return children;
};
