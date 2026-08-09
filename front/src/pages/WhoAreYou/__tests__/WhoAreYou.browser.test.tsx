import { appRoutes } from "@Front/routing/appRoutes";
import { renderBrowserRoute } from "@Front/utils/testsUtils/customRender/customRender.browser";
import { page, userEvent } from "vitest/browser";

describe("WhoAreYou Page", () => {
  it("should render page with all fields", async () => {
    await renderBrowserRoute({ initialEntry: appRoutes.whoAreYou() });

    await expect
      .element(page.getByRole("heading", { level: 1, name: "Who are you?" }))
      .toBeInTheDocument();
    await expect.element(page.getByLabelText(/Avatar/u)).toBeInTheDocument();
    await expect.element(page.getByLabelText(/Username/u)).toBeInTheDocument();
    await expect
      .element(page.getByLabelText(/Favorite color/u))
      .toBeInTheDocument();
    await expect
      .element(page.getByLabelText(/I accept the terms and conditions of use/u))
      .toBeInTheDocument();
    await expect
      .element(page.getByRole("button", { name: "Continuer" }))
      .toBeInTheDocument();
  });

  it("should show errors when submitting empty form", async () => {
    await renderBrowserRoute({ initialEntry: appRoutes.whoAreYou() });

    await page.getByRole("button", { name: "Continuer" }).click();

    await expect
      .element(page.getByText("The avatar is required"))
      .toBeInTheDocument();
    await expect
      .element(page.getByText("Username is required"))
      .toBeInTheDocument();
    await expect
      .element(page.getByText("You must choose a color"))
      .toBeInTheDocument();
    await expect
      .element(
        page.getByText("You must accept the terms and conditions of use"),
      )
      .toBeInTheDocument();
  });

  describe("Avatar field validation", () => {
    it("should show error when uploading an invalid avatar", async () => {
      await renderBrowserRoute({ initialEntry: appRoutes.whoAreYou() });

      const fileInput = page.getByLabelText(/Avatar/u);
      const invalidFile = new File(["invalid content"], "invalid.txt", {
        type: "text/plain",
      });

      await userEvent.upload(fileInput, invalidFile);
      await page.getByRole("button", { name: "Continuer" }).click();

      await expect
        .element(
          page.getByText("The avatar file must be a PNG, JPEG, or WEBP image"),
        )
        .toBeInTheDocument();
    });

    it("should show error when uploading an avatar that is too large", async () => {
      await renderBrowserRoute({ initialEntry: appRoutes.whoAreYou() });

      const fileInput = page.getByLabelText(/Avatar/u);
      const largeFile = new File(
        [new ArrayBuffer(11 * 1024 * 1024)],
        "large.png",
        {
          type: "image/png",
        },
      );

      await userEvent.upload(fileInput, largeFile);
      await page.getByRole("button", { name: "Continuer" }).click();

      await expect
        .element(page.getByText("The avatar file size must be less than 10 MB"))
        .toBeInTheDocument();
    });
  });

  describe("Username field validation", () => {
    it("should show error when username is too short", async () => {
      await renderBrowserRoute({ initialEntry: appRoutes.whoAreYou() });

      const usernameInput = page.getByLabelText(/Username/u);
      await userEvent.fill(usernameInput, "ab");
      await page.getByRole("button", { name: "Continuer" }).click();

      await expect
        .element(page.getByText("Username must be at least 3 characters"))
        .toBeInTheDocument();
    });

    it("should show error when username is too long", async () => {
      await renderBrowserRoute({ initialEntry: appRoutes.whoAreYou() });

      const usernameInput = page.getByLabelText(/Username/u);
      await userEvent.fill(usernameInput, "a".repeat(31));
      await page.getByRole("button", { name: "Continuer" }).click();

      await expect
        .element(page.getByText("Username must be at most 30 characters"))
        .toBeInTheDocument();
    });
  });
});
