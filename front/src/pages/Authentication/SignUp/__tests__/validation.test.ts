import type { TFunction } from "i18next";
import { describe, expect, it } from "vitest";
import { getSchema } from "../validation";

const translate = ((key) => String(key)) as TFunction;
const schema = getSchema(translate);
const validSignUp = {
  email: "user@example.com",
  password: "Password1!",
  confirmPassword: "Password1!",
  termsAccepted: true,
};

describe("SignUp validation schema", () => {
  it("should accept valid sign-up data", async () => {
    await expect(schema.validate(validSignUp)).resolves.toEqual(validSignUp);
  });

  describe("email", () => {
    it("should require an email", async () => {
      await expect(
        schema.validateAt("email", { ...validSignUp, email: "" }),
      ).rejects.toThrow("requiredEmail");
    });

    it.each(["not-an-email", "user@domain", "user @example.com"])(
      "should reject an invalid email: %s",
      async (email) => {
        await expect(
          schema.validateAt("email", { ...validSignUp, email }),
        ).rejects.toThrow("invalidEmail");
      },
    );
  });

  describe("password", () => {
    it("should require a password", async () => {
      await expect(
        schema.validateAt("password", { ...validSignUp, password: "" }),
      ).rejects.toThrow("requiredPassword");
    });

    it("should enforce the minimum length", async () => {
      await expect(
        schema.validateAt("password", {
          ...validSignUp,
          password: "Short1!",
        }),
      ).rejects.toThrow("minLengthPassword");
    });

    it("should enforce password complexity", async () => {
      await expect(
        schema.validateAt("password", {
          ...validSignUp,
          password: "password1!",
        }),
      ).rejects.toThrow("passwordComplexity");
    });
  });

  describe("password confirmation", () => {
    it("should require a confirmation", async () => {
      await expect(
        schema.validateAt("confirmPassword", {
          ...validSignUp,
          password: "",
          confirmPassword: "",
        }),
      ).rejects.toThrow("requiredConfirmPassword");
    });

    it("should match the password", async () => {
      await expect(
        schema.validateAt("confirmPassword", {
          ...validSignUp,
          confirmPassword: "Different1!",
        }),
      ).rejects.toThrow("passwordsDoNotMatch");
    });
  });

  describe("terms acceptance", () => {
    it("should require the terms to be accepted", async () => {
      await expect(
        schema.validateAt("termsAccepted", {
          ...validSignUp,
          termsAccepted: false,
        }),
      ).rejects.toThrow("termsAcceptedError");
    });
  });
});
