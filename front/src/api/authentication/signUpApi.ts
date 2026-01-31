import type {
  SignUpErrorCodeType,
  SignUpRequestBodyType,
  SignUpResponseType,
} from '@Front/types/Authentication/signUp/signUp.types';
import { SignUpErrorResponse } from '@Front/types/Authentication/signUp/SignUpErrorResponse';
import { METHODS } from '../constant';
import { fetchApi } from '../fetchApi';

export const signUpApi = ({
  username,
  email,
  password,
  language,
}: SignUpRequestBodyType): Promise<SignUpResponseType> => {
  const body: SignUpRequestBodyType = { username, email, password, language };

  return fetchApi<SignUpResponseType, SignUpErrorCodeType>({
    path: `${import.meta.env.FRONT_BACKEND_URL}/v1/account`,
    method: METHODS.post,
    data: body,
    CustomErrorResponse: SignUpErrorResponse,
  });
};
