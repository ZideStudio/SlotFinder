import { appRoutes } from "@Front/routing/appRoutes";
import { renderBrowserRoute } from "@Front/utils/testsUtils/customRender/customRender.browser";
import { worker } from "@Mocks/browser";
import {
  patchAccount200,
  patchAccount400,
  patchAvatarAccount200,
  patchAvatarAccount400,
} from "@Mocks/handlers/accountHandlers";
import { page } from "vitest/browser";

describe("WhoAreYou Page", () => {
  const fillFormWithValidData = async () => {
    const avatarInput = page.getByLabelText(/Avatar/u);
    const validAvatar = new File(["avatar"], "avatar.png", {
      type: "image/png",
    });

    await avatarInput.upload(validAvatar);
    await page.getByLabelText(/Username/u).fill("john_doe");
    await page.getByLabelText(/Favorite color/u).fill("#ff0000");
    await page
      .getByRole("checkbox", {
        name: /I accept the terms and conditions of use/u,
      })
      .click();
  };

  afterEach(() => {
    worker.resetHandlers();
  });

  it("should render page with all fields", async () => {
    await renderBrowserRoute({ initialEntry: appRoutes.whoAreYou() });

    await expect
      .element(page.getByRole("heading", { level: 1, name: "Who are you ?" }))
      .toBeInTheDocument();
    await expect.element(page.getByLabelText(/Avatar/u)).toBeInTheDocument();
    await expect
      .element(page.getByRole("textbox", { name: /Username/u }))
      .toBeInTheDocument();
    await expect
      .element(page.getByLabelText(/Favorite color/u))
      .toBeInTheDocument();
    await expect
      .element(
        page.getByRole("checkbox", {
          name: /I accept the terms and conditions of use/u,
        }),
      )
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

      const invalidFile = new File(["invalid content"], "invalid.txt", {
        type: "text/plain",
      });

      await page.getByLabelText(/Avatar/u).upload(invalidFile);
      await page.getByRole("button", { name: "Continuer" }).click();

      await expect
        .element(
          page.getByText("The avatar file must be a PNG, JPEG, or WEBP image"),
        )
        .toBeInTheDocument();
    });

    it("should show error when uploading an avatar that is too large", async () => {
      await renderBrowserRoute({ initialEntry: appRoutes.whoAreYou() });

      const largeFile = new File(
        [new ArrayBuffer(11 * 1024 * 1024)],
        "large.png",
        {
          type: "image/png",
        },
      );

      await page.getByLabelText(/Avatar/u).upload(largeFile);
      await page.getByRole("button", { name: "Continuer" }).click();

      await expect
        .element(page.getByText("The avatar file size must be less than 10 MB"))
        .toBeInTheDocument();
    });
  });

  describe("Username field validation", () => {
    it("should show error when username is too short", async () => {
      await renderBrowserRoute({ initialEntry: appRoutes.whoAreYou() });

      await page.getByRole("textbox", { name: "Username" }).fill("ab");
      await page.getByRole("button", { name: "Continuer" }).click();

      await expect
        .element(page.getByText("Username must be at least 3 characters"))
        .toBeInTheDocument();
    });

    it("should show error when username is too long", async () => {
      await renderBrowserRoute({ initialEntry: appRoutes.whoAreYou() });

      await page
        .getByRole("textbox", { name: "Username" })
        .fill("a".repeat(31));
      await page.getByRole("button", { name: "Continuer" }).click();

      await expect
        .element(page.getByText("Username must be at most 30 characters"))
        .toBeInTheDocument();
    });
  });

  describe("Submit API errors", () => {
    it("should show username error when account patch returns username already taken", async () => {
      worker.use(
        patchAccount400("USERNAME_ALREADY_TAKEN"),
        patchAvatarAccount200,
      );

      await renderBrowserRoute({ initialEntry: appRoutes.whoAreYou() });
      await fillFormWithValidData();

      await page.getByRole("button", { name: "Continuer" }).click();

      await expect
        .element(
          page.getByText(
            "This username is already taken. Please choose another one.",
          ),
        )
        .toBeInTheDocument();
    });

    it("should show color error when account patch returns invalid color format", async () => {
      worker.use(
        patchAccount400("INVALID_COLOR_FORMAT"),
        patchAvatarAccount200,
      );

      await renderBrowserRoute({ initialEntry: appRoutes.whoAreYou() });
      await fillFormWithValidData();

      await page.getByRole("button", { name: "Continuer" }).click();

      await expect
        .element(
          page.getByText(
            "The color must be a valid hexadecimal code (ex: #RRGGBB).",
          ),
        )
        .toBeInTheDocument();
    });

    it("should show server error when account patch returns server error", async () => {
      worker.use(patchAccount400("SERVER_ERROR"), patchAvatarAccount200);

      await renderBrowserRoute({ initialEntry: appRoutes.whoAreYou() });
      await fillFormWithValidData();

      await page.getByRole("button", { name: "Continuer" }).click();

      await expect
        .element(
          page.getByText(
            "An unexpected error occurred during submit. Please try again later or contact support if the issue persists.",
          ),
        )
        .toBeInTheDocument();
    });

    it("should show avatar error when avatar upload fails", async () => {
      worker.use(patchAccount200, patchAvatarAccount400);

      await renderBrowserRoute({ initialEntry: appRoutes.whoAreYou() });
      await fillFormWithValidData();

      await page.getByRole("button", { name: "Continuer" }).click();

      await expect
        .element(
          page.getByText("Failed to upload the avatar. Please try again."),
        )
        .toBeInTheDocument();
    });
  });
});
