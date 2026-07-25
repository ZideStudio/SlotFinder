import { ErrorResponse } from "@Front/types/ErrorResponse";
import { HEADERS } from "../constant";
import { tokenRefreshManager } from "../tokenRefreshManager";

/**
 * Custom fetch mutator used by Orval-generated API clients.
 *
 * Receives the path from the prepared OpenAPI spec (e.g. `/v1/account`),
 * which works as a relative URL in the browser. The body and headers are
 * already prepared by the generated code (JSON.stringify, Content-Type, …).
 *
 * Returns the parsed response body (or `undefined` for empty bodies) to match
 * the types Orval generates for the `fetch` client when
 * `includeHttpResponseReturnType: false`. Throws an `ErrorResponse` on non-2xx
 * responses so React Query captures it as `mutation.error` and callers can
 * still use `getErrorCode()` to read the structured error payload.
 *
 * Replicates the core behaviour of `fetchApi`:
 * - Handles 498 (expired access token) by refreshing and retrying once.
 * - Forwards the AbortSignal provided by React Query / TanStack Query.
 */
export const fetchApiMutator = async <Response>(
  url: string,
  options: RequestInit,
  signal?: AbortSignal,
): Promise<Response> => {
  const apiUrlFull = `${window.location.origin}${import.meta.env.FRONT_BACKEND_URL ?? ""}`;

  const makeRequest = async (): Promise<globalThis.Response> =>
    await fetch(url, { ...options, signal });

  let response = await makeRequest();

  // Handle 498 status code from api (expired access token)
  if (response.status === 498 && response.url.startsWith(apiUrlFull)) {
    await tokenRefreshManager.refreshToken();
    // Retry the original request
    response = await makeRequest();
  }

  const content = await response.text();

  if (!response.ok) {
    throw new ErrorResponse(content);
  }

  // Some endpoints legitimately return an empty body (e.g. 204 or void responses)
  if (!content) {
    return undefined as Response;
  }

  if ((response.headers.get(HEADERS.contentType) ?? "").includes("json")) {
    return JSON.parse(content) as Response;
  }

  return content as Response;
};
