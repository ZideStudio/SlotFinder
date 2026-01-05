import type { Json } from '@Front/types/api.types';
import { ErrorResponse } from '@Front/types/ErrorResponse';
import { isAuthenticatedPage } from '@Front/utils/isAuthenticatedPage';
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

// Singleton for managing token refresh to prevent multiple simultaneous refresh calls
class TokenRefreshManager {
  private static instance: TokenRefreshManager;
  private isRefreshing = false;
  private refreshPromise: Promise<void> | null = null;

  private constructor() {}

  public static getInstance(): TokenRefreshManager {
    if (!TokenRefreshManager.instance) {
      TokenRefreshManager.instance = new TokenRefreshManager();
    }
    return TokenRefreshManager.instance;
  }

  public async refreshToken(): Promise<void> {
    // If already refreshing, wait for that operation to complete
    if (this.isRefreshing && this.refreshPromise) {
      await this.refreshPromise;
      return;
    }

    // Start a new refresh operation
    this.isRefreshing = true;
    this.refreshPromise = this.performRefresh();

    try {
      await this.refreshPromise;
    } finally {
      this.isRefreshing = false;
      this.refreshPromise = null;
    }
  }

  private async performRefresh(): Promise<void> {
    const response = await fetch(`${import.meta.env.FRONT_BACKEND_URL}/v1/auth/refresh`, {
      method: 'POST',
      credentials: 'include',
    });

    if (!response.ok) {
      // On refresh failure, redirect to home page
      window.location.href = '/';
      throw new Error('Token refresh failed');
    }
  }
}

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
  // Only attempt refresh if we're on an authenticated page and it's not the refresh endpoint itself
  if (response.status === 498) {
    const isOnAuthenticatedPage = isAuthenticatedPage();
    const isRefreshEndpoint = path.includes('/auth/refresh');

    if (isOnAuthenticatedPage && !isRefreshEndpoint) {
      try {
        const refreshManager = TokenRefreshManager.getInstance();
        await refreshManager.refreshToken();
        
        // Retry the original request after successful refresh
        response = await makeRequest();
      } catch (error) {
        // If refresh failed, the manager already redirected to home
        // Just rethrow to prevent further processing
        throw error;
      }
    }
  }

  // Handle 401 status code (unauthorized) - do nothing, just let it pass through
  // This indicates the refresh token has expired or is invalid

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
