import { appRoutes } from "@Front/routing/appRoutes";
import { renderBrowserRoute } from "@Front/utils/testsUtils/customRender/customRender.browser";
import { worker } from "@Mocks/browser";
import {
  getAccountMe200,
  getAccountMeWithoutTerms200,
  patchAccount200,
  patchAvatarAccount200,
  patchAvatarAccount400,
} from "@Mocks/handlers/accountHandlers";
import { getAuthStatus403 } from "@Mocks/handlers/authStatusHandlers";
import { page } from "vitest/browser";

/**
 * Browser smoke tests for WhoAreYou page.
 *
 * These tests focus on real browser rendering, accessibility, and core user flows.
 * Detailed form validation and API error handling are covered in unit tests.
 * See: useWhoAreYou.test.ts for comprehensive business logic coverage.
 */
describe("WhoAreYou Page", () => {
  const fillFormWithValidData = async () => {
    const avatarInput = page.getByLabelText(/Avatar/u);
    const validAvatar = new File(["avatar"], "avatar.png", {
      type: "image/png",
    });

    await avatarInput.upload(validAvatar);
    await page.getByLabelText(/Username/u).fill("john_doe");
    await page.getByLabelText(/Favorite color/u).fill("#ff0000");
  };

  afterEach(() => {
    worker.resetHandlers();
  });

  it("should hide terms checkbox when terms are already accepted at current version", async () => {
    worker.use(getAuthStatus403("TERMS_NOT_ACCEPTED"), getAccountMe200);

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
      .not.toBeInTheDocument();
    await expect
      .element(page.getByRole("button", { name: "Continue" }))
      .toBeInTheDocument();
  });

  it("should show terms checkbox when terms are not accepted", async () => {
    worker.use(
      getAuthStatus403("TERMS_NOT_ACCEPTED"),
      getAccountMeWithoutTerms200,
    );

    await renderBrowserRoute({ initialEntry: appRoutes.whoAreYou() });

    await expect
      .element(
        page.getByRole("checkbox", {
          name: /I accept the terms and conditions of use/u,
        }),
      )
      .toBeInTheDocument();
  });

  it("should display error message when avatar upload fails", async () => {
    worker.use(
      getAuthStatus403("USERNAME_MISSING"),
      getAccountMe200,
      patchAccount200,
      patchAvatarAccount400,
    );

    await renderBrowserRoute({ initialEntry: appRoutes.whoAreYou() });
    await fillFormWithValidData();

    await page.getByRole("button", { name: "Continue" }).click();

    await expect
      .element(page.getByText("Failed to upload the avatar. Please try again."))
      .toBeInTheDocument();
  });

  it("should submit valid data successfully without error message", async () => {
    worker.use(
      getAuthStatus403("USERNAME_MISSING"),
      getAccountMe200,
      patchAccount200,
      patchAvatarAccount200,
    );

    await renderBrowserRoute({ initialEntry: appRoutes.whoAreYou() });
    await fillFormWithValidData();

    await page.getByRole("button", { name: "Continue" }).click();

    await expect
      .element(
        page.getByText(
          "An unexpected error occurred during submit. Please try again later or contact support if the issue persists.",
        ),
      )
      .not.toBeInTheDocument();
  });
});
