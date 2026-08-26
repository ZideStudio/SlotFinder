import { page } from "@vitest/browser/context";
import { render } from "vitest-browser-react";
import LoaderPage from "../LoaderPage";

describe("LoaderPage", () => {
  beforeEach(() => {
    vi.useFakeTimers();
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  it("does not show the loader immediately on mount", async () => {
    render(<LoaderPage />);

    await expect
      .element(page.getByRole("status", { name: "loading" }))
      .not.toBeInTheDocument();
  });

  it("does not show the loader just before 100ms", async () => {
    render(<LoaderPage />);

    await vi.advanceTimersByTimeAsync(99);

    await expect
      .element(page.getByRole("status", { name: "loading" }))
      .not.toBeInTheDocument();
  });

  it("shows the loader after 100ms", async () => {
    render(<LoaderPage />);

    await vi.advanceTimersByTimeAsync(100);

    await expect
      .element(page.getByRole("status", { name: "loading" }))
      .toBeInTheDocument();
  });
});