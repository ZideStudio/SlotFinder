import type { Json } from '@Front/types/api.types';
import { ErrorResponse } from '@Front/types/ErrorResponse';
import { HEADERS, METHODS, MIME_TYPES } from './constant';

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

export const fetchApi = async <Response extends Json | string, CustomErrorResponseCodeType extends string = never>({
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

  const response = await fetch(path, {
    method,
    ...(data && { body: JSON.stringify(data) }),
    headers: mergeHeaders,
  });

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
