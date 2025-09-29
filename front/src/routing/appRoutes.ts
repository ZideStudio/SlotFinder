import { signUpRoutes } from '@Front/pages/Authentication/SignUp';

export const appRoutes: Record<string, () => string> = {
  home: () => '/',
  signUp: () => `/${signUpRoutes.path}`,
};
