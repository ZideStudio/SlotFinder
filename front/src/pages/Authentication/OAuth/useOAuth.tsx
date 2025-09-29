import { oAuthProvidersApi } from '@Front/api/authentication/oAuthProvidersApi';
import type {
  OAuthProvidersErrorCodeType,
  OAuthProvidersResponseType,
} from '@Front/types/Authentication/oAuthProviders/oAuthProviders.types';
import type { OAuthProvidersErrorResponse } from '@Front/types/Authentication/oAuthProviders/OAuthProvidersErrorResponse';
import { useQuery } from '@tanstack/react-query';
import { useMemo } from 'react';
import { oauthProvidersData, ONE_DAY_IN_MS } from './constants';
import type { OAuthProvider } from './types';

type TUseOAuth = {
  oAuthProviders: OAuthProvider[];
  errorCode?: OAuthProvidersErrorCodeType;
};

export const useOAuth = (): TUseOAuth => {
  const { data, error } = useQuery<OAuthProvidersResponseType, OAuthProvidersErrorResponse>({
    queryKey: ['oAuthProviders'],
    queryFn: oAuthProvidersApi,
    staleTime: ONE_DAY_IN_MS,
    gcTime: ONE_DAY_IN_MS,
  });

  const oAuthProviders: OAuthProvider[] = useMemo(
    () =>
      oauthProvidersData.map(provider => ({
        ...provider,
        href: data ? data[provider.id] : '#',
      })),
    [data],
  );

  const errorCode = useMemo(() => error?.getErrorCode(), [error]);

  return {
    oAuthProviders,
    errorCode,
  };
};
