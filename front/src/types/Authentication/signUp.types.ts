export type SignUpRequestBodyType = {
  username: string;
  email: string;
  password: string;
};

export type SignUpFormType = SignUpRequestBodyType;

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

export type SignUpErrorType = {
  code: string;
  error: true;
  message: string;
};
