import { userEvent } from "@vitest/browser/context";
import { render } from "vitest-browser-react";
import { page } from "vitest/browser";

import { OverlayContent } from "../OverlayContent";

describe("OverlayContent", () => {
  describe("rendering", () => {
    it("should render title as a level 1 heading, content, and footer actions", async () => {
      await render(
        <OverlayContent
          title="Content title"
          primaryButtonProps={{ children: "Confirm" }}
          secondaryButtonProps={{ children: "Cancel" }}
        >
          Body content
        </OverlayContent>,
      );

      await expect
        .element(page.getByRole("heading", { name: "Content title", level: 1 }))
        .toBeInTheDocument();
      await expect.element(page.getByText("Body content")).toBeInTheDocument();
      await expect
        .element(page.getByRole("button", { name: "Confirm" }))
        .toBeInTheDocument();
      await expect
        .element(page.getByRole("button", { name: "Cancel" }))
        .toBeInTheDocument();
    });

    it("should not render secondary button when secondaryButtonProps is not provided", async () => {
      await render(
        <OverlayContent
          title="Content title"
          primaryButtonProps={{ children: "Confirm" }}
        >
          Body content
        </OverlayContent>,
      );

      await expect
        .element(page.getByRole("button", { name: "Confirm" }))
        .toBeInTheDocument();
      await expect
        .element(page.getByRole("button", { name: "Cancel" }))
        .not.toBeInTheDocument();
    });

    it("should render a close button with accessible name and structural sections", async () => {
      await render(
        <OverlayContent
          title="Content title"
          primaryButtonProps={{ children: "Confirm" }}
        >
          Body content
        </OverlayContent>,
      );

      const closeButton = page.getByRole("button", {
        name: "Fermer la fenêtre",
      });
      await expect.element(closeButton).toBeInTheDocument();
      await expect
        .element(closeButton)
        .toHaveAccessibleName("Fermer la fenêtre");

      await expect.element(page.getByRole("banner")).toBeInTheDocument();
      await expect.element(page.getByText("Body content")).toBeInTheDocument();

      const footer = page.getByRole("contentinfo");
      await expect.element(footer).toBeInTheDocument();
      await expect
        .element(footer.getByRole("button", { name: "Confirm" }))
        .toBeInTheDocument();
    });
  });

  describe("titleId", () => {
    it("should set the id attribute on the heading when titleId is provided", async () => {
      await render(
        <OverlayContent
          title="Content title"
          titleId="custom-title-id"
          primaryButtonProps={{ children: "Confirm" }}
        >
          Body content
        </OverlayContent>,
      );

      await expect
        .element(page.getByRole("heading", { name: "Content title" }))
        .toHaveAttribute("id", "custom-title-id");
    });

    it("should not set the id attribute on the heading when titleId is not provided", async () => {
      await render(
        <OverlayContent
          title="Content title"
          primaryButtonProps={{ children: "Confirm" }}
        >
          Body content
        </OverlayContent>,
      );

      await expect
        .element(page.getByRole("heading", { name: "Content title" }))
        .not.toHaveAttribute("id");
    });
  });

  describe("interactions", () => {
    it("should call closeButtonProps.onClick when the close button is clicked", async () => {
      const onClick = vi.fn();

      await render(
        <OverlayContent
          title="Content title"
          primaryButtonProps={{ children: "Confirm" }}
          closeButtonProps={{ onClick }}
        >
          Body content
        </OverlayContent>,
      );

      await userEvent.click(
        page.getByRole("button", { name: "Fermer la fenêtre" }),
      );
      expect(onClick).toHaveBeenCalledTimes(1);
    });

    it("should call primaryButtonProps onClick when primary button is clicked", async () => {
      const onClick = vi.fn();

      await render(
        <OverlayContent
          title="Content title"
          primaryButtonProps={{ children: "Confirm", onClick }}
        >
          Body content
        </OverlayContent>,
      );

      await userEvent.click(page.getByRole("button", { name: "Confirm" }));
      expect(onClick).toHaveBeenCalledTimes(1);
    });

    it("should call secondaryButtonProps onClick when secondary button is clicked", async () => {
      const onClick = vi.fn();

      await render(
        <OverlayContent
          title="Content title"
          primaryButtonProps={{ children: "Confirm" }}
          secondaryButtonProps={{ children: "Cancel", onClick }}
        >
          Body content
        </OverlayContent>,
      );

      await userEvent.click(page.getByRole("button", { name: "Cancel" }));
      expect(onClick).toHaveBeenCalledTimes(1);
    });
  });
});
