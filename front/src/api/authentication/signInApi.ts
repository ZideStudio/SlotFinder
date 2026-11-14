import {
  type SigninTokenResponseDto,
  type SigninSigninDto,
} from "@Front/api/generated/slotFinderAPI.schemas";
import type { SignInRequestBodyType } from "@Front/types/Authentication/signIn/signIn.types";
import { postV1AuthSignin } from "../generated";

export const signInApi = ({
  identifier,
  password,
}: SignInRequestBodyType): Promise<SigninTokenResponseDto> => {
  const dto: SigninSigninDto = {
    identifier,
    password,
  };

  return postV1AuthSignin(dto);
};
