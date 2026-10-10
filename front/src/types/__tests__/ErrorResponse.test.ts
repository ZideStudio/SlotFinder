import { describe, expect, it } from "vitest";
import { ErrorResponse } from "@Front/types/ErrorResponse";
import {
  getErrorCodeFromPayload,
  getErrorMessageFromPayload,
} from "@Front/types/ErrorResponse.helpers";

describe("ErrorResponse helpers", () => {
  it("should extract code from a JSON payload", () => {
    expect(getErrorCodeFromPayload('{"code":"USERNAME_ALREADY_TAKEN"}')).toBe(
      "USERNAME_ALREADY_TAKEN",
    );
  });

  it("should extract code from a nested error payload", () => {
    expect(
      getErrorCodeFromPayload({
        error: { code: "INVALID_COLOR_FORMAT", message: "bad color" },
      }),
    ).toBe("INVALID_COLOR_FORMAT");
  });

  it("should extract the message from a payload", () => {
    expect(
      getErrorMessageFromPayload({
        message: "Request failed",
      }),
    ).toBe("Request failed");
  });
});

describe("ErrorResponse", () => {
  it("should normalize a JSON payload into an ErrorResponse instance", () => {
    const error = ErrorResponse.from({
      code: "USERNAME_ALREADY_TAKEN",
      message: "Username already taken",
    });

    expect(error).toBeInstanceOf(ErrorResponse);
    expect(error.getErrorCode()).toBe("USERNAME_ALREADY_TAKEN");
    expect(error.message).toBe("Username already taken");
  });

  it("should fall back to SERVER_ERROR when the payload is not structured", () => {
    const error = ErrorResponse.from(new Error("boom"));

    expect(error.getErrorCode()).toBe("SERVER_ERROR");
    expect(error.message).toBe("boom");
  });
});
