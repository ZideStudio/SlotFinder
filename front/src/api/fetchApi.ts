import type { Json } from '@Front/types/api.types';
import { ErrorResponse } from '@Front/types/ErrorResponse';
import { HEADERS, METHODS, MIME_TYPES } from './constant';
import { tokenRefreshManager } from './tokenRefreshManager';

type FetchApiError = Error & {
  code?: number;
  contentType?: MIME_TYPES;
};

type ErrorResponseClass<ErrorCodeType extends string> = new (message: string) => ErrorResponse<ErrorCodeType>;

type FetchApiProps<CustomErrorResponseCodeType extends string> = {
  path: string;
  method?: METHODS;
  data?: Json;
  headers?: HeadersInit;
  CustomErrorResponse?: ErrorResponseClass<CustomErrorResponseCodeType>;
};

export const fetchApi = async <
  Response extends Json | string | null,
  CustomErrorResponseCodeType extends string = never,
>({
  path,
  method = METHODS.get,
  data,
  headers = [],
  CustomErrorResponse = ErrorResponse,
}: FetchApiProps<CustomErrorResponseCodeType>): Promise<Response> => {
  const mergeHeaders = new Headers(headers);

  if (data) {
    mergeHeaders.append(HEADERS.contentType, MIME_TYPES.json);
  }

  const makeRequest = async (): Promise<globalThis.Response> => {
    return await fetch(path, {
      method,
      ...(data && { body: JSON.stringify(data) }),
      headers: mergeHeaders,
      credentials: 'include',
    });
  };

  let response = await makeRequest();

  // Handle 498 status code (expired access token)
  if (response.status === 498) {
    await tokenRefreshManager.refreshToken();
    response = await makeRequest(); // Retry the original request
  }

  const content = await response.text();

  if (!response.ok) {
    const error: FetchApiError = new CustomErrorResponse(content);
    error.code = response.status;
    error.contentType = (response.headers.get(HEADERS.contentType) as MIME_TYPES) ?? MIME_TYPES.text;

    throw error;
  }

  if ((response.headers.get(HEADERS.contentType) ?? '').includes('json')) {
    return JSON.parse(content) as Response;
  }

  return content as Response;
};
