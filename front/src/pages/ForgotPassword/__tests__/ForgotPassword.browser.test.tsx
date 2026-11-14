import { appRoutes } from "@Front/routing/appRoutes";
import {
  renderRoute,
  type RenderRouteOptions,
} from "@Front/utils/testsUtils/customRender/customRender";
import {
  postForgotPassword200,
  postForgotPassword400,
} from "@Mocks/handlers/accountHandlers";
import { worker } from "@Mocks/browser";
import { screen } from "@testing-library/react";
import { userEvent } from "@testing-library/user-event";
import { authenticationRoutes } from "../../Authentication/routes";

const renderRouteOptions: RenderRouteOptions = {
  routes: [authenticationRoutes],
  routesOptions: { initialEntries: [appRoutes.forgotPassword()] },
};

describe("ForgotPassword", () => {
  beforeEach(() => {
    worker.use(postForgotPassword200);
  });

  afterEach(() => {
    worker.resetHandlers();
  });

  it("renders the email field and submit button", () => {
    renderRoute(renderRouteOptions);

    expect(
      screen.queryByRole("navigation", {
        name: "Sign in with a provider",
      }),
    ).not.toBeInTheDocument();
    expect(
      screen.getByRole("heading", { name: "Forgot your password?" }),
    ).toBeInTheDocument();
    expect(screen.getByRole("textbox", { name: /^Email/u })).toBeRequired();
    expect(
      screen.getByRole("button", { name: "Send reset link" }),
    ).toBeInTheDocument();
  });

  it("shows accessible validation errors for an empty or invalid email", async () => {
    renderRoute(renderRouteOptions);

    await userEvent.click(
      screen.getByRole("button", { name: "Send reset link" }),
    );

    const emailInput = screen.getByRole("textbox", { name: /^Email/u });
    const requiredError = await screen.findByText("Email is required");
    expect(requiredError).toBeVisible();
    expect(emailInput).toHaveAttribute("aria-invalid", "true");
    expect(emailInput).toHaveAttribute("aria-describedby", requiredError.id);

    await userEvent.type(emailInput, "not-an-email");
    await userEvent.click(
      screen.getByRole("button", { name: "Send reset link" }),
    );

    const invalidEmailError = await screen.findByText(
      "Enter a valid email address",
    );
    expect(invalidEmailError).toBeVisible();
    expect(emailInput).toHaveAttribute(
      "aria-describedby",
      invalidEmailError.id,
    );
  });

  it("shows an explicit error when reset requests are rate limited", async () => {
    worker.use(postForgotPassword400);
    renderRoute(renderRouteOptions);
    await userEvent.type(
      screen.getByRole("textbox", { name: /^Email/u }),
      "person@example.com",
    );
    await userEvent.click(
      screen.getByRole("button", { name: "Send reset link" }),
    );

    await expect(
      screen.findByText(
        "A reset request was sent recently. Please wait before trying again.",
      ),
    ).resolves.toBeVisible();
  });
});
