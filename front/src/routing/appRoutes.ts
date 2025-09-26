import { signUpRoutes } from '@Front/pages/SignUp';

export const appRoutes: Record<string, () => string> = {
  home: () => '/',
  signUp: () => `/${signUpRoutes.path}`,
};
