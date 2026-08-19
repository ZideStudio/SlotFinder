import { userEvent } from "@vitest/browser/context";
import { render } from "vitest-browser-react";
import { page } from "vitest/browser";

import { Modal } from "../Modal";

describe("Modal", () => {
  it("should render title, content, and footer actions", async () => {
    await render(
      <Modal
        open
        title="Modal title"
        primaryButtonProps={{ children: "Action" }}
        secondaryButtonProps={{ children: "Close" }}
      >
        Modal content
      </Modal>,
    );

    const dialog = page.getByRole("dialog");
    await expect.element(dialog).toBeInTheDocument();
    await expect
      .element(dialog.getByRole("heading", { name: "Modal title" }))
      .toBeInTheDocument();
    await expect.element(dialog.getByText("Modal content")).toBeInTheDocument();
    await expect
      .element(dialog.getByRole("button", { name: "Action" }))
      .toBeInTheDocument();
    await expect
      .element(dialog.getByRole("button", { name: "Close" }))
      .toBeInTheDocument();
  });

  it("should link dialog label to the title heading and set closedby attribute", async () => {
    await render(
      <Modal
        open
        title="Modal title"
        primaryButtonProps={{ children: "Action" }}
      >
        Modal content
      </Modal>,
    );

    const dialog = page.getByRole("dialog");
    const heading = page.getByRole("heading", { name: "Modal title" });

    await expect.element(dialog).toBeInTheDocument();
    await expect
      .element(dialog)
      .toHaveAttribute("aria-labelledby", (await heading.element()).id);
    await expect.element(dialog).toHaveAttribute("closedby", "any");
  });

  it("should apply the variant passed via secondaryButtonProps", async () => {
    await render(
      <Modal
        open
        title="Modal title"
        primaryButtonProps={{ children: "Action" }}
        secondaryButtonProps={{ children: "Close", variant: "secondary" }}
      >
        Modal content
      </Modal>,
    );

    await expect
      .element(page.getByRole("button", { name: "Close" }))
      .toHaveClass("ds-button--secondary");
  });

  it("should forward native dialog attributes", async () => {
    const onClose = vi.fn();
    await render(
      <Modal
        title="Modal title"
        primaryButtonProps={{ children: "Action" }}
        open
        onClose={onClose}
      >
        Modal content
      </Modal>,
    );

    await expect.element(page.getByRole("dialog")).toHaveAttribute("open");
  });

  it("should apply a custom className alongside the default class", async () => {
    await render(
      <Modal
        open
        title="Modal title"
        primaryButtonProps={{ children: "Action" }}
        className="custom-modal"
      >
        Modal content
      </Modal>,
    );

    const dialog = page.getByRole("dialog");
    await expect.element(dialog).toHaveClass("ds-modal");
    await expect.element(dialog).toHaveClass("custom-modal");
  });

  it("should call closeModal when the close button is clicked", async () => {
    const close = vi.fn();
    Object.defineProperty(HTMLDialogElement.prototype, "close", {
      value: close,
      configurable: true,
      writable: true,
    });

    await render(
      <Modal
        open
        title="Modal title"
        primaryButtonProps={{ children: "Action" }}
      >
        Modal content
      </Modal>,
    );

    await userEvent.click(
      page.getByRole("button", { name: "Fermer la fenêtre" }),
    );
    expect(close).toHaveBeenCalledTimes(1);
  });
});
