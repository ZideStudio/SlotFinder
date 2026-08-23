import { useAuthenticationContext } from "@Front/hooks/useAuthenticationContext";
import { appRoutes } from "@Front/routing/appRoutes";
import type { RouteHandle } from "@Front/routing/routeHandle";
import type { AuthStatusErrorCodeType } from "@Front/types/Authentication/authStatus/authStatus.types";
import { useMemo, type ReactNode } from "react";
import { Navigate, useLocation, useMatches, type UIMatch } from "react-router";

const REDIRECT_TO_WHO_ARE_YOU_ERRORS = new Set<AuthStatusErrorCodeType>([
  "USERNAME_MISSING",
  "TERMS_NOT_ACCEPTED",
]);

type AuthenticationProtectionProps = {
  children: ReactNode;
};

export const AuthenticationProtection = ({
  children,
}: AuthenticationProtectionProps) => {
  const {
    isAuthenticated,
    authenticationError,
    postAuthRedirectPath,
    setPostAuthRedirectPath,
    resetPostAuthRedirectPath,
  } = useAuthenticationContext();
  const { pathname } = useLocation();
  const matches = useMatches() as UIMatch<unknown, RouteHandle>[];

  const mustBeAuthenticate = useMemo(() => {
    const currentMatch = matches.at(-1);

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

  if (
    authenticationError &&
    REDIRECT_TO_WHO_ARE_YOU_ERRORS.has(authenticationError) &&
    pathname !== appRoutes.whoAreYou()
  ) {
    return <Navigate to={appRoutes.whoAreYou()} replace />;
  }

  if (
    mustBeAuthenticate &&
    authenticationError &&
    !REDIRECT_TO_WHO_ARE_YOU_ERRORS.has(authenticationError)
  ) {
    return <Navigate to={appRoutes.signUp()} replace />;
  }

  if (mustBeAuthenticate === false && isAuthenticated) {
    return <Navigate to={postAuthRedirectPath ?? appRoutes.home()} replace />;
  }

  return children;
};
