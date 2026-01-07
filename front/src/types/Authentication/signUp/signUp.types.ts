import type { ErrorResponseCodeType } from '@Front/types/api.types';

export type SignUpRequestBodyType = {
  username: string;
  email: string;
  password: string;
  language: string;
};

export type SignUpFormType = Omit<SignUpRequestBodyType, "language"> & {
  confirmPassword: string;
};

export type SignUpResponseType = {
  access_token: string;
  createdAt: string;
  email: string;
  id: string;
  providers:
    | [
        {
          provider: string;
        },
      ]
    | null;
  userName: string;
};

export type SignUpErrorCodeType = ErrorResponseCodeType<
  'USERNAME_ALREADY_TAKEN' | 'INVALID_EMAIL_FORMAT' | 'EMAIL_ALREADY_EXISTS' | 'INVALID_PASSWORD_FORMAT'
>;
