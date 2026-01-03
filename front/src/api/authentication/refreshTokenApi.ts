import { METHODS } from '../constant';
import { fetchApi } from '../fetchApi';

export const refreshTokenApi = () =>
  fetchApi<null>({
    path: `${import.meta.env.FRONT_BACKEND_URL}/v1/auth/refresh`,
    method: METHODS.post,
  });
