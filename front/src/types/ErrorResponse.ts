import { SERVER_ERROR } from "@Front/utils/constants/api";

export class ErrorResponse<ErrorCodeType extends string = never> extends Error {
  constructor(message: string) {
    super(message);
    this.name = "ErrorResponse";
  }

  getErrorCode(
    this: ErrorResponse<ErrorCodeType>,
  ): ErrorCodeType | "SERVER_ERROR";
  getErrorCode(this: ErrorResponse<never>): "SERVER_ERROR" {
    try {
      const parsed = JSON.parse(this.message);
      if (
        parsed &&
        typeof parsed === "object" &&
        "code" in parsed &&
        typeof parsed.code === "string"
      ) {
        return parsed.code;
      }
    } catch {
      // Ignore
    }
    return SERVER_ERROR;
  }
}
