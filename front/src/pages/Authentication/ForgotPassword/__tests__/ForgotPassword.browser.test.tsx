import { appRoutes } from "@Front/routing/appRoutes";
import {
  renderRoute,
  type RenderRouteOptions,
} from "@Front/utils/testsUtils/customRender/customRender";
import { postForgotPassword200 } from "@Mocks/handlers/accountHandlers";
import { worker } from "@Mocks/browser";
import { screen } from "@testing-library/react";
import { userEvent } from "@testing-library/user-event";
import { authenticationRoutes } from "../../../Authentication/routes";
import { afterEach, beforeEach, describe, expect, it } from "vitest";

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

  it("shows an accessible validation error for an invalid email", async () => {
    renderRoute(renderRouteOptions);

    const emailInput = screen.getByRole("textbox", { name: /^Email/u });
    await userEvent.type(emailInput, "not-an-email");
    await userEvent.click(
      screen.getByRole("button", { name: "Send reset link" }),
    );

    const invalidEmailError = await screen.findByText(
      "Enter a valid email address",
    );
    expect(invalidEmailError).toBeVisible();
    expect(emailInput).toHaveAttribute("aria-invalid", "true");
    expect(emailInput).toHaveAttribute(
      "aria-describedby",
      invalidEmailError.id,
    );
  });

  it("shows a confirmation message when the reset request is successful", async () => {
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
        "If an account exists for this email address, you’ll receive a password reset link shortly.",
      ),
    ).resolves.toBeVisible();
  });
});
