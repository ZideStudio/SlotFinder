import type { AuthStatusErrorCodeType } from "@Front/types/Authentication/authStatus/authStatus.types";
import type { UseCheckAuthenticationReturn } from "./useCheckAuthentication";
import type { UsePostAuthenticationReturn } from "./usePostAuthentication";

export type AuthenticationContextType = {
  isAuthenticated: boolean | undefined;
  authenticationError: AuthStatusErrorCodeType | undefined;
  checkAuthentication: UseCheckAuthenticationReturn["checkAuthentication"];
  postAuthRedirectPath: UsePostAuthenticationReturn["postAuthRedirectPath"];
  setPostAuthRedirectPath: UsePostAuthenticationReturn["setPostAuthRedirectPath"];
  resetPostAuthRedirectPath: UsePostAuthenticationReturn["resetPostAuthRedirectPath"];
};
