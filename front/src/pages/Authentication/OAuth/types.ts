import type { OAuthProvidersResponseType } from '@Front/types/Authentication/oAuthProviders/oAuthProviders.types';
import type { ReactNode } from 'react';

export type OAuthProvider = {
  id: keyof OAuthProvidersResponseType;
  label: string;
  href: string;
  ariaLabel: `signInWith${Capitalize<keyof OAuthProvidersResponseType>}`;
  icon: ReactNode;
};
