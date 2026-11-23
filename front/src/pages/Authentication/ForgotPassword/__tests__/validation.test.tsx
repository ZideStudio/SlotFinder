import { type TFunction } from "i18next";
import { getSchema } from "../validation";

const translate = ((key) => String(key)) as TFunction<"forgotPassword">;
const schema = getSchema(translate);
const validForgotPassword = {
  email: "",
};

describe("email", () => {
  it("should require an email", async () => {
    await expect(
      schema.validateAt("email", { ...validForgotPassword, email: "" }),
    ).rejects.toThrow("requiredEmail");
  });

  it.each(["not-an-email", "user@domain", "user @example.com"])(
    "should reject an invalid email: %s",
    async (email) => {
      await expect(
        schema.validateAt("email", { ...validForgotPassword, email }),
      ).rejects.toThrow("invalidEmail");
    },
  );

  it("should accept a valid email", async () => {
    await expect(
      schema.validateAt("email", {
        ...validForgotPassword,
        email: "person@example.com",
      }),
    ).resolves.toBe("person@example.com");
  });

  it("should accept an email with 320 characters", async () => {
    const email = `${"a".repeat(308)}@example.com`;

    await expect(email).toHaveLength(320);
    await expect(
      schema.validateAt("email", { ...validForgotPassword, email }),
    ).resolves.toBe(email);
  });

  it("should reject an email with 321 characters", async () => {
    const email = `${"a".repeat(309)}@example.com`;

    await expect(email).toHaveLength(321);
    await expect(
      schema.validateAt("email", { ...validForgotPassword, email }),
    ).rejects.toThrow("invalidEmail");
  });
});
