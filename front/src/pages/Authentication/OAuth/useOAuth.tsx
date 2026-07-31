import { getGetV1AuthProviderUrlUrl } from "@Front/api/generated/authentication/authentication";
import { useAuthenticationContext } from "@Front/hooks/useAuthenticationContext";
import { appRoutes } from "@Front/routing/appRoutes";
import { oauthProvidersData } from "./constants";
import type { OAuthProvider } from "./types";

type TUseOAuth = {
  oAuthProviders: OAuthProvider[];
};

export const useOAuth = (): TUseOAuth => {
  const { postAuthRedirectPath } = useAuthenticationContext();
  const oAuthProviders: OAuthProvider[] = oauthProvidersData.map(
    (provider) => ({
      ...provider,
      href: getGetV1AuthProviderUrlUrl(provider.id, {
        returnUrl: postAuthRedirectPath || appRoutes.home(),
      }),
    }),
  );

  return {
    oAuthProviders,
  };
};
