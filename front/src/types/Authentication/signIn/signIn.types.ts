export type SignInRequestBodyType = {
  identifier: string;
  password: string;
};

export type SignInResponseType = {
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
