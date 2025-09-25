import type {
  SignUpFormType,
  SignUpRequestBodyType,
  SignUpResponseType,
} from '@Front/types/Authentication/signUp.types';
import { METHODS } from './constant';
import { fetchApi } from './fetchApi';

export const signUpApi = async ({ username, email, password }: SignUpFormType): Promise<SignUpResponseType> => {
  const body: SignUpRequestBodyType = { username, email, password };

  return await fetchApi<SignUpResponseType>({
    path: `${import.meta.env.FRONT_BACKEND_URL}/v1/account`,
    method: METHODS.post,
    data: body,
  });
};
