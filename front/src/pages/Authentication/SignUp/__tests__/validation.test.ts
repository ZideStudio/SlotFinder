import type { TFunction } from "i18next";
import { describe, expect, it } from "vitest";
import { getSchema } from "../validation";

const translate = ((key) => String(key)) as TFunction<"signUp">;
const schema = getSchema(translate);
const validSignUp = {
  email: "user@example.com",
  password: "Password1!",
  confirmPassword: "Password1!",
  termsAccepted: true,
};

describe("SignUp validation schema", () => {
  it("should accept valid sign-up data", () => {
    expect(schema.parse(validSignUp)).toStrictEqual(validSignUp);
  });

  describe("email", () => {
    it("should require an email", () => {
      expect(() => schema.parse({ ...validSignUp, email: "" })).toThrow(
        "requiredEmail",
      );
    });

    it.each(["not-an-email", "user@domain", "user @example.com"])(
      "should reject an invalid email: %s",
      (email) => {
        expect(() => schema.parse({ ...validSignUp, email })).toThrow(
          "invalidEmail",
        );
      },
    );
  });

  describe("password", () => {
    it("should require a password", () => {
      expect(() => schema.parse({ ...validSignUp, password: "" })).toThrow(
        "requiredPassword",
      );
    });

    it("should enforce the minimum length", () => {
      expect(() =>
        schema.parse({
          ...validSignUp,
          password: "Short1!",
        }),
      ).toThrow("minLengthPassword");
    });

    it("should enforce password complexity", () => {
      expect(() =>
        schema.parse({
          ...validSignUp,
          password: "password1!",
        }),
      ).toThrow("passwordComplexity");
    });
  });

  describe("password confirmation", () => {
    it("should require a confirmation", () => {
      expect(() =>
        schema.parse({
          ...validSignUp,
          password: "",
          confirmPassword: "",
        }),
      ).toThrow("requiredConfirmPassword");
    });

    it("should match the password", () => {
      expect(() =>
        schema.parse({
          ...validSignUp,
          confirmPassword: "Different1!",
        }),
      ).toThrow("passwordsDoNotMatch");
    });
  });

  describe("terms acceptance", () => {
    it("should require the terms to be accepted", () => {
      expect(() =>
        schema.parse({
          ...validSignUp,
          termsAccepted: false,
        }),
      ).toThrow("termsAcceptedError");
    });
  });
});
