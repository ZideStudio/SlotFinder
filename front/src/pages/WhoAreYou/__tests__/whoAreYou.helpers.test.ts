
import { TERMS_VERSION } from "@Front/utils/constants/terms";
import type { WhoAreYouFormData } from "../types";
import { buildAccountUpdateDTO } from "../whoAreYou.helpers";

describe("whoAreYou helpers", () => {
  const originalTemporal = globalThis.Temporal;

  beforeAll(() => {
    globalThis.Temporal = {
      Now: {
        timeZoneId: () => "Europe/Paris",
      },
    } as unknown as typeof Temporal;
  });

  afterAll(() => {
    globalThis.Temporal = originalTemporal;
  });

  describe("buildAccountUpdateDTO", () => {
    it("should build correct DTO from form data", () => {
      const formData: WhoAreYouFormData = {
        username: "john_doe",
        color: "#ff0000",
        termsAccepted: true,
        avatar: [] as unknown as FileList,
      };

      const result = buildAccountUpdateDTO(formData);

      expect(result).toStrictEqual({
        username: "john_doe",
        color: "#ff0000",
        termsAccepted: true,
        termsVersion: TERMS_VERSION,
        timeZone: "Europe/Paris",
      });
    });

    it("should include timeZone from Temporal API", () => {
      const formData: WhoAreYouFormData = {
        username: "test_user",
        color: "#00ff00",
        termsAccepted: false,
        avatar: [] as unknown as FileList,
      };

      const result = buildAccountUpdateDTO(formData);

      expect(result.timeZone).toBe("Europe/Paris");
    });
  });
});
