import { server } from "@Mocks/server";
import "@testing-library/jest-dom/vitest";
import { cleanup } from "@testing-library/react";
import * as matchers from "vitest-axe/matchers";

// Polyfill DataTransfer which jsdom does not implement
class MockDataTransfer {
  private fileArray: File[] = [];

  get items() {
    return {
      add: (file: File) => {
        this.fileArray.push(file);
      },
    };
  }

  get files(): FileList {
    const result = [...this.fileArray] as unknown as FileList;
    Object.assign(result, {
      item: (index: number) => this.fileArray[index] || null,
    });
    return result;
  }
}

// Polyfill DataTransfer for Node.js
if (typeof globalThis.DataTransfer === "undefined") {
  (globalThis as unknown as Record<string, unknown>).DataTransfer =
    MockDataTransfer;
}

expect.extend(matchers);

beforeAll(() => {
  server.listen({ onUnhandledRequest: "error" });

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
});

afterAll(() => {
  server.close();
});

afterEach(() => {
  server.resetHandlers();
  cleanup();
});
