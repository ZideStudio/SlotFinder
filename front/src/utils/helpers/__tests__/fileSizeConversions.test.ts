import { bytesToMB, mbToBytes } from "../fileSizeConversions";

describe("fileSizeConversions", () => {
  describe("bytesToMB", () => {
    it("should convert bytes to megabytes with two decimal places", () => {
      expect(bytesToMB(1048576)).toBe(1);
      expect(bytesToMB(1572864)).toBe(1.5);
    });

    it("should round the converted value to two decimal places", () => {
      expect(bytesToMB(1536000)).toBe(1.46);
    });

    it("should return 0 for 0 bytes", () => {
      expect(bytesToMB(0)).toBe(0);
    });
  });

  describe("mbToBytes", () => {
    it("should convert megabytes to bytes", () => {
      expect(mbToBytes(1)).toBe(1048576);
      expect(mbToBytes(1.5)).toBe(1572864);
    });

    it("should return 0 for 0 megabytes", () => {
      expect(mbToBytes(0)).toBe(0);
    });
  });
});
