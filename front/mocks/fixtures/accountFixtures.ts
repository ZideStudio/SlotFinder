import type { SignUpErrorType, SignUpResponseType } from "@Front/types/Authentication/signUp.types";

export const account: SignUpResponseType = {
  access_token: '1234567890abcdef',
  createdAt: '2024-01-01T00:00:00.000Z',
  email: 'test@example.com',
  id: '123456',
  providers: null,
  userName: 'test_user',
};

export const accountError: SignUpErrorType = {
  code: 'TEST_ERROR',
  error: true,
  message: 'This is a test error message on account creation.',
};
