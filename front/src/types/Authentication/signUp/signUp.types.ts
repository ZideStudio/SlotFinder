import type { ErrorResponseCodeType } from "@Front/types/api.types";

export type SignUpRequestBodyType = {
  email: string;
  password: string;
  language: string;
  termsAccepted: boolean;
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
};

export type SignUpErrorCodeType = ErrorResponseCodeType<
  | "INVALID_EMAIL_FORMAT"
  | "EMAIL_ALREADY_EXISTS"
  | "INVALID_PASSWORD_FORMAT"
>;
