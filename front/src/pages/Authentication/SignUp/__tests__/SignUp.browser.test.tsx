import { appRoutes } from "@Front/routing/appRoutes";
import { renderBrowserRoute } from "@Front/utils/testsUtils/customRender/customRender.browser";
import { worker } from "@Mocks/browser";
import { getAuthStatus401 } from "@Mocks/handlers/authStatusHandlers";
import { postAccount201 } from "@Mocks/handlers/accountHandlers";
import { page } from "vitest/browser";

describe("SignUp page", () => {
  beforeEach(() => {
    worker.use(getAuthStatus401);
  });

  afterEach(() => {
    worker.resetHandlers();
  });

  it("renders the sign-up form with all expected inputs", async () => {
    worker.use(postAccount201);

    await renderBrowserRoute({ initialEntry: appRoutes.signUp() });

    await expect
      .element(
        page.getByRole("heading", {
          level: 1,
          name: "Sign Up to SlotFinder",
        }),
      )
      .toBeInTheDocument();
    await expect.element(page.getByLabelText("Email")).toBeInTheDocument();
    await expect.element(page.getByLabelText("Password")).toBeInTheDocument();
    await expect
      .element(page.getByLabelText("Confirm password"))
      .toBeInTheDocument();
    await expect
      .element(
        page.getByRole("checkbox", {
          name: "I agree to the Terms and Conditions",
        }),
      )
      .toBeInTheDocument();
    await expect
      .element(page.getByRole("button", { name: "Create Account" }))
      .toBeInTheDocument();
  });

  it("shows validation errors when the required fields are empty", async () => {
    await renderBrowserRoute({ initialEntry: appRoutes.signUp() });

    await page.getByRole("button", { name: "Create Account" }).click();

    await expect
      .element(page.getByText("Email is required"))
      .toBeInTheDocument();
    await expect
      .element(page.getByText("Password is required"))
      .toBeInTheDocument();
    await expect
      .element(page.getByText("Confirm password is required"))
      .toBeInTheDocument();
    await expect
      .element(page.getByText("You must accept the Terms and Conditions"))
      .toBeInTheDocument();
  });

  it("shows a mismatch error when the password confirmation is different", async () => {
    await renderBrowserRoute({ initialEntry: appRoutes.signUp() });

    await page.getByLabelText("Email").fill("john@example.com");
    await page.getByLabelText("Password").fill("Password1!");
    await page.getByLabelText("Confirm password").fill("Password2!");
    await page
      .getByRole("checkbox", { name: "I agree to the Terms and Conditions" })
      .click();

    await page.getByRole("button", { name: "Create Account" }).click();

    await expect
      .element(page.getByText("Passwords do not match"))
      .toBeInTheDocument();
  });
});
