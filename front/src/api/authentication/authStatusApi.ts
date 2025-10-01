import type { AuthStatusErrorCodeType, AuthStatusResponseType } from "@Front/types/Authentication/authStatus/authStatus.types";
import { AuthStatusErrorResponse } from "@Front/types/Authentication/authStatus/AuthStatusErrorResponse";
import { METHODS } from "../constant";
import { fetchApi } from "../fetchApi";

export const authStatusApi = () =>
  fetchApi<AuthStatusResponseType, AuthStatusErrorCodeType>({
    path: `${import.meta.env.FRONT_BACKEND_URL}/v1/auth/status`,
    method: METHODS.get,
    sendCredentials: true,
    CustomErrorResponse: AuthStatusErrorResponse,
  });