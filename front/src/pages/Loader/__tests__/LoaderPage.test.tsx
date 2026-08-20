import React from "react";
import { render, screen, act } from "@testing-library/react";
import LoaderPage from "../LoaderPage";

describe("LoaderPage", () => {
  beforeEach(() => {
    vi.useFakeTimers();
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  it("show nothing immediately on mount", () => {
    render(<LoaderPage />);

    expect(screen.queryByRole('status', { name: 'loading' })).not.toBeInTheDocument();
  });

  it("show nothing just before 100ms", async () => {
    render(<LoaderPage />);

    await vi.advanceTimersByTimeAsync(99);

    expect(screen.queryByRole('status', { name: 'loading' })).not.toBeInTheDocument();
  });

  it("show the loader after 100ms", async () => {
    render(<LoaderPage />);

    await vi.advanceTimersByTimeAsync(100);

    expect(screen.getByRole('status', { name: 'loading' })).toBeInTheDocument();
  });

  it("never shows the loader if it is unmounted before 100ms", async () => {
    const { unmount } = render(<LoaderPage />);

    await vi.advanceTimersByTime(50);
    unmount();
    await vi.advanceTimersByTime(100);

    expect(screen.queryByRole('status', { name: 'loading' })).not.toBeInTheDocument();
  });
});