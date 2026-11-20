import type { TFunction } from "i18next";
import { describe, expect, it } from "vitest";
import { AVATAR_MAX_SIZE_MB } from "@Front/pages/WhoAreYou/constants";
import { getUserSchema } from "../../validation/userSchema";

const translate = ((key) => String(key)) as TFunction;
const schema = getUserSchema(translate);

const createFileList = (...files: File[]): FileList => {
  const fileList = Object.create(FileList.prototype) as FileList;

  files.forEach((file, index) => {
    Object.defineProperty(fileList, index, { value: file });
  });

  Object.defineProperties(fileList, {
    length: { value: files.length },
    item: { value: (index: number) => files[index] ?? null },
  });

  return fileList;
};

const validData = {
  email: "user@example.com",
  avatar: createFileList(
    new File(["avatar"], "avatar.png", { type: "image/png" }),
  ),
  username: "sam",
  color: "#12aBcF",
  termsAccepted: true,
  password: "Password1!",
  confirmPassword: "Password1!",
};

describe("getCommonSchema", () => {
  it("should accept valid data", async () => {
    await expect(schema.validate(validData)).resolves.toEqual(validData);
  });

  describe("email", () => {
    it("should reject an invalid email", async () => {
      await expect(
        schema.validateAt("email", { ...validData, email: "invalid-email" }),
      ).rejects.toThrow("invalidEmail");
    });

    it("should reject an email longer than 320 characters", async () => {
      const email = `${"a".repeat(309)}@example.com`;

      await expect(
        schema.validateAt("email", { ...validData, email }),
      ).rejects.toThrow("invalidEmail");
    });
  });

  describe("avatar", () => {
    it("should require a non-empty file list", async () => {
      await expect(
        schema.validateAt("avatar", { ...validData, avatar: undefined }),
      ).rejects.toThrow("avatarRequired");

      await expect(
        schema.validateAt("avatar", {
          ...validData,
          avatar: createFileList(),
        }),
      ).rejects.toThrow("avatarRequired");
    });

    it("should reject files that exceed the maximum size", async () => {
      const oversizedFile = new File(
        [new Uint8Array(AVATAR_MAX_SIZE_MB * 1024 * 1024 + 1)],
        "avatar.png",
        { type: "image/png" },
      );

      await expect(
        schema.validateAt("avatar", {
          ...validData,
          avatar: createFileList(oversizedFile),
        }),
      ).rejects.toThrow("avatarFileSizeError");
    });

    it("should reject unsupported file types", async () => {
      const unsupportedFile = new File(["avatar"], "avatar.gif", {
        type: "image/gif",
      });

      await expect(
        schema.validateAt("avatar", {
          ...validData,
          avatar: createFileList(unsupportedFile),
        }),
      ).rejects.toThrow("avatarFileTypeError");
    });
  });

  describe("username", () => {
    it("should require a username", async () => {
      await expect(
        schema.validateAt("username", { ...validData, username: "" }),
      ).rejects.toThrow("requiredUsername");
    });

    it("should enforce the minimum length", async () => {
      await expect(
        schema.validateAt("username", { ...validData, username: "ab" }),
      ).rejects.toThrow("minLengthUsername");
    });

    it("should enforce the maximum length", async () => {
      await expect(
        schema.validateAt("username", {
          ...validData,
          username: "a".repeat(31),
        }),
      ).rejects.toThrow("maxLengthUsername");
    });
  });

  describe("color", () => {
    it("should require a color", async () => {
      await expect(
        schema.validateAt("color", { ...validData, color: "" }),
      ).rejects.toThrow("colorRequired");
    });

    it("should reject an invalid color format", async () => {
      await expect(
        schema.validateAt("color", { ...validData, color: "red" }),
      ).rejects.toThrow("colorInvalid");
    });
  });

  describe("termsAccepted", () => {
    it("should require the terms to be accepted", async () => {
      await expect(
        schema.validateAt("termsAccepted", {
          ...validData,
          termsAccepted: false,
        }),
      ).rejects.toThrow("termsAcceptedError");
    });
  });
});
