import { useToastService } from "@Front/ui/utils/toast/hooks/useToastService";
import { userEvent } from "@vitest/browser/context";
import { render } from "vitest-browser-react";
import { page } from "vitest/browser";
import { ToastProvider } from "../ToastProvider";

describe("ToastProvider", () => {
  afterEach(() => {
    vi.useRealTimers();
  });

  const TestComponent = () => {
    const toast = useToastService();
    return (
      <button onClick={() => toast.addToast("Test Toast")}>Show Toast</button>
    );
  };

  it("should render children and provide toast context", async () => {
    await render(
      <ToastProvider>
        <TestComponent />
      </ToastProvider>,
    );

    const button = page.getByRole("button", { name: "Show Toast" });
    await expect.element(button).toBeInTheDocument();

    await userEvent.click(button);

    await expect.element(page.getByText("Test Toast")).toBeInTheDocument();
  });

  it("should remove toast after duration", async () => {
    vi.useFakeTimers();

    await render(
      <ToastProvider>
        <TestComponent />
      </ToastProvider>,
    );

    await userEvent.click(page.getByRole("button", { name: "Show Toast" }));

    await expect.element(page.getByText("Test Toast")).toBeInTheDocument();

    await vi.runAllTimersAsync();

    await expect.element(page.getByText("Test Toast")).not.toBeInTheDocument();
  });

  it("should remove toast when close button is clicked", async () => {
    await render(
      <ToastProvider>
        <TestComponent />
      </ToastProvider>,
    );

    await userEvent.click(page.getByText("Show Toast"));

    await expect.element(page.getByText("Test Toast")).toBeInTheDocument();

    await userEvent.click(
      page.getByRole("button", { name: "Fermer la notification" }),
    );
    await expect.element(page.getByText("Test Toast")).not.toBeInTheDocument();
  });
});
