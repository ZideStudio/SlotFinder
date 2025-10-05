import type { UseCheckAuthenticationReturn } from './useCheckAuthentication';

export type AuthenticationContextType = {
  isAuthenticated: boolean | undefined;
  checkAuthentication: UseCheckAuthenticationReturn['checkAuthentication'];
};
