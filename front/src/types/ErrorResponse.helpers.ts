import { SERVER_ERROR } from "@Front/utils/constants/api";

const UNEXPECTED_ERROR_MESSAGE = "An unexpected error occurred.";

const isObject = (value: unknown): value is Record<string, unknown> =>
  typeof value === "object" && value !== null;

export const getErrorCodeFromPayload = (payload: unknown): string => {
  if (typeof payload === "string") {
    try {
      return getErrorCodeFromPayload(JSON.parse(payload));
    } catch {
      return SERVER_ERROR;
    }
  }

  if (!isObject(payload)) {
    return SERVER_ERROR;
  }

  const value = payload as Record<string, unknown>;

  if (typeof value.code === "string") {
    return value.code;
  }

  if (value.error && isObject(value.error)) {
    const nested = value.error as Record<string, unknown>;
    if (typeof nested.code === "string") {
      return nested.code;
    }
  }

  return SERVER_ERROR;
};

export const getErrorMessageFromPayload = (payload: unknown): string => {
  if (typeof payload === "string") {
    return payload;
  }

  if (payload instanceof Error) {
    return payload.message;
  }

  if (!isObject(payload)) {
    return UNEXPECTED_ERROR_MESSAGE;
  }

  const value = payload as Record<string, unknown>;

  if (typeof value.message === "string") {
    return value.message;
  }

  if (value.error && isObject(value.error)) {
    const nested = value.error as Record<string, unknown>;
    if (typeof nested.message === "string") {
      return nested.message;
    }
  }

  return UNEXPECTED_ERROR_MESSAGE;
};
