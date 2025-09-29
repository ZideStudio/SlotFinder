import type { ErrorResponseType } from '@Front/types/api.types';
import { type SignUpErrorCodeType, type SignUpResponseType } from '@Front/types/Authentication/signUp/signUp.types';

export const accountFixture: SignUpResponseType = {
  access_token: '1234567890abcdef',
  createdAt: '2024-01-01T00:00:00.000Z',
  email: 'test@example.com',
  id: '123456',
  providers: null,
  userName: 'test_user',
};

export const accountErrorFixture: ErrorResponseType<SignUpErrorCodeType> = {
  code: 'USERNAME_ALREADY_TAKEN',
};
