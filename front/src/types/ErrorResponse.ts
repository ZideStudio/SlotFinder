import { SERVER_ERROR } from "@Front/utils/constants/api";
import {
  getErrorCodeFromPayload,
  getErrorMessageFromPayload,
} from "./ErrorResponse.helpers";

export class ErrorResponse<ErrorCodeType extends string = never> extends Error {
  code: ErrorCodeType | "SERVER_ERROR";

  static from(payload: unknown): ErrorResponse<string> {
    if (payload instanceof ErrorResponse) {
      return payload as ErrorResponse<string>;
    }

    const message = getErrorMessageFromPayload(payload);
    const code = getErrorCodeFromPayload(payload) as string;

    return new ErrorResponse<string>(message, code as "SERVER_ERROR" | string);
  }

  constructor(
    message: string,
    code: ErrorCodeType | "SERVER_ERROR" = SERVER_ERROR,
  ) {
    super(message);
    this.name = "ErrorResponse";
    this.code = code;
  }

  getErrorCode(): ErrorCodeType | "SERVER_ERROR" {
    if (this.code) {
      return this.code;
    }

    const fallbackCode = getErrorCodeFromPayload(this.message) as
      | ErrorCodeType
      | "SERVER_ERROR";
    this.code = fallbackCode;

    return fallbackCode;
  }
}
