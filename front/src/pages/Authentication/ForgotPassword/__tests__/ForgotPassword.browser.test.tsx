import { appRoutes } from "@Front/routing/appRoutes";
import {
  renderRoute,
  type RenderRouteOptions,
} from "@Front/utils/testsUtils/customRender/customRender";
import {
  postForgotPassword200,
  postForgotPassword400,
  postForgotPassword500,
} from "@Mocks/handlers/accountHandlers";
import { worker } from "@Mocks/browser";
import { screen } from "@testing-library/react";
import { userEvent } from "@testing-library/user-event";
import { authenticationRoutes } from "../../../Authentication/routes";

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

  it("shows an error message when the reset request fails", async () => {
    worker.use(postForgotPassword400);
    renderRoute(renderRouteOptions);

    await userEvent.type(
      screen.getByRole("textbox", { name: /^Email/u }),
      "person@example.com",
    );
    await userEvent.click(
      screen.getByRole("button", { name: "Send reset link" }),
    );

    const errorMessage = await screen.findByRole("alert");
    expect(errorMessage).toBeVisible();
    expect(errorMessage).toHaveTextContent(
      "A reset request was sent recently. Please wait before trying again.",
    );
  });

  it("shows SERVER_ERROR when the reset request fails with an unknown error", async () => {
    worker.use(postForgotPassword500);
    renderRoute(renderRouteOptions);

    await userEvent.type(
      screen.getByRole("textbox", { name: /^Email/u }),
      "person@example.com",
    );
    await userEvent.click(
      screen.getByRole("button", { name: "Send reset link" }),
    );

    const errorMessage = await screen.findByRole("alert");
    expect(errorMessage).toBeVisible();
    expect(errorMessage).toHaveTextContent(
      "We couldn’t send the reset link. Please try again later.",
    );
  });
});
