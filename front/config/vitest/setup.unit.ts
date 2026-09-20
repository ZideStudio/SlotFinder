import { server } from "@Mocks/server";
import "@testing-library/jest-dom/vitest";
import { cleanup } from "@testing-library/react";
import * as matchers from "vitest-axe/matchers";

vi.mock("react-i18next", () => ({
  useTranslation: vi.fn((resource: string) => ({
    t: (messageId: string, args: Record<string, unknown>) =>
      `${resource}.${messageId}${args ? `::${JSON.stringify(args)}` : ""}`,
    i18n: {
      language: "en",
    },
  })),
  initReactI18next: {
    type: "3rdParty",
    init: () => {},
  },
}));

expect.extend(matchers);

beforeAll(() => {
  server.listen({ onUnhandledRequest: "error" });
});

afterAll(() => {
  server.close();
});

afterEach(() => {
  server.resetHandlers();
  cleanup();
});
