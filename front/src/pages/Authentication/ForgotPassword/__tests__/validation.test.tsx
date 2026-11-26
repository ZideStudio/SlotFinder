import { describe, expect, it } from "vitest";
import { type TFunction } from "i18next";
import { getSchema } from "../validation";

const translate = ((key) => String(key)) as TFunction<"forgotPassword">;
const schema = getSchema(translate);
const validForgotPassword = {
  email: "person@example.com",
};

describe("email", () => {
  it("should require an email", () => {
    expect(() => schema.parse({ ...validForgotPassword, email: "" })).toThrow(
      "requiredEmail",
    );
  });

  it.each(["not-an-email", "user@domain", "user @example.com"])(
    "should reject an invalid email: %s",
    (email) => {
      expect(() => schema.parse({ ...validForgotPassword, email })).toThrow(
        "invalidEmail",
      );
    },
  );

  it("should accept a valid email", () => {
    expect(
      schema.parse({
        ...validForgotPassword,
        email: "person@example.com",
      }).email,
    ).toBe("person@example.com");
  });

  it("should accept an email with 320 characters", async () => {
    const email = `${"a".repeat(308)}@example.com`;

    await expect(email).toHaveLength(320);
    expect(schema.parse({ ...validForgotPassword, email }).email).toBe(email);
  });

  it("should reject an email with 321 characters", async () => {
    const email = `${"a".repeat(309)}@example.com`;

    await expect(email).toHaveLength(321);
    expect(() => schema.parse({ ...validForgotPassword, email })).toThrow(
      "invalidEmail",
    );
  });
});
