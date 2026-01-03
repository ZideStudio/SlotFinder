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
  skipAuthRefresh?: boolean; // Flag to skip automatic token refresh
};

// Global state for managing token refresh
let isRefreshing = false;
let refreshPromise: Promise<void> | null = null;

const refreshToken = async (): Promise<void> => {
  const response = await fetch(`${import.meta.env.FRONT_BACKEND_URL}/v1/auth/refresh`, {
    method: 'POST',
    credentials: 'include',
  });

  if (!response.ok) {
    throw new Error('Token refresh failed');
  }
};

const waitForTokenRefresh = async (): Promise<void> => {
  if (refreshPromise) {
    await refreshPromise;
  }
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
  skipAuthRefresh = false,
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

  // Handle 401 error with automatic token refresh
  if (response.status === 401 && !skipAuthRefresh) {
    // Check if path is an authenticated endpoint
    const isAuthenticatedEndpoint = !path.includes('/auth/status') && !path.includes('/auth/refresh');

    if (isAuthenticatedEndpoint) {
      try {
        // If another request is already refreshing, wait for it
        if (isRefreshing) {
          await waitForTokenRefresh();
          // Retry the original request after refresh completes
          response = await makeRequest();
        } else {
          // Start the refresh process
          isRefreshing = true;
          refreshPromise = refreshToken()
            .then(() => {
              isRefreshing = false;
              refreshPromise = null;
            })
            .catch((error) => {
              isRefreshing = false;
              refreshPromise = null;
              // Redirect to home on refresh failure
              window.location.href = '/';
              throw error;
            });

          await refreshPromise;

          // Retry the original request
          response = await makeRequest();
        }
      } catch (error) {
        // If refresh failed, the catch block already redirected
        // Just rethrow to prevent further processing
        throw error;
      }
    }
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
