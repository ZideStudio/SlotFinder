import { appRoutes } from '@Front/routing/appRoutes';
import { oauthProvidersData } from './constants';
import type { OAuthProvider, OAuthProviderName } from './types';

type TUseOAuth = {
  oAuthProviders: OAuthProvider[];
};

/**
 * Generates the OAuth URL for a given provider and redirect URL.
 *
 * @param provider - The OAuth provider name (e.g., 'discord', 'google', 'github').
 * @param redirectUrl - The URL to which the user will be redirected after authentication.
 * @returns The complete OAuth authorization URL for the specified provider.
 */
const generateOAuthUrl = (provider: OAuthProviderName, redirectUrl: string): string =>
  `${import.meta.env.FRONT_BACKEND_URL}/v1/auth/${provider}/url?redirectUrl=${encodeURIComponent(redirectUrl)}`;

export const useOAuth = (): TUseOAuth => {
  const oAuthProviders: OAuthProvider[] = oauthProvidersData.map(provider => ({
    ...provider,
    href: generateOAuthUrl(provider.id, appRoutes.dashboard()),
  }));

  return {
    oAuthProviders,
  };
};
