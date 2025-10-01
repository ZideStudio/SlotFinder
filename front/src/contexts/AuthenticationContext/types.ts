import type { UseCheckAuthenticationReturn } from './useCheckAuthentication';

export type AuthenticationContextType = {
  isAuthenticated: boolean;
  checkAuthentication: UseCheckAuthenticationReturn['checkAuthentication'];
};
