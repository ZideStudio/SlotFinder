export type JsonArray = unknown[];
export type JsonObject = Record<string, unknown>;
export type Json = JsonObject | JsonArray;

export type ErrorResponseCodeType<OtherError extends string = never> = 'SERVER_ERROR' | OtherError;
export type ErrorResponseType<ErrorCodeType extends string = never> = {
  code: ErrorCodeType;
};
