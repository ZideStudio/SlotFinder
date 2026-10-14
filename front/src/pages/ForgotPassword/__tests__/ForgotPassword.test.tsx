import { appRoutes } from "@Front/routing/appRoutes";
import {
  renderRoute,
  type RenderRouteOptions,
} from "@Front/utils/testsUtils/customRender/customRender";
import {
  postForgotPassword200,
  postForgotPassword400,
} from "@Mocks/handlers/accountHandlers";
import { server } from "@Mocks/server";
import { screen } from "@testing-library/react";
import { userEvent } from "@testing-library/user-event";
import { authenticationRoutes } from "../../Authentication/routes";

const renderRouteOptions: RenderRouteOptions = {
  routes: [authenticationRoutes],
  routesOptions: { initialEntries: [appRoutes.forgotPassword()] },
};

describe("ForgotPassword", () => {
  beforeEach(() => {
    server.use(postForgotPassword200);
  });

  afterEach(() => {
    server.resetHandlers();
  });

  it("renders the email field and submit button", () => {
    renderRoute(renderRouteOptions);

    expect(
      screen.getByRole("heading", { name: "forgotPassword.title" }),
    ).toBeInTheDocument();
    expect(screen.getByLabelText(/forgotPassword\.email/u)).toBeRequired();
    expect(
      screen.getByRole("button", { name: "forgotPassword.submit" }),
    ).toBeInTheDocument();
  });

  it("shows accessible validation errors for an empty or invalid email", async () => {
    renderRoute(renderRouteOptions);

    await userEvent.click(
      screen.getByRole("button", { name: "forgotPassword.submit" }),
    );

    const emailInput = screen.getByLabelText(/forgotPassword\.email/u);
    await expect(
      screen.findByText("forgotPassword.requiredEmail"),
    ).resolves.toBeVisible();
    expect(emailInput).toHaveAttribute("aria-invalid", "true");
    expect(emailInput).toHaveAttribute(
      "aria-describedby",
      "forgot-password-email-error",
    );

    await userEvent.type(emailInput, "not-an-email");
    await userEvent.click(
      screen.getByRole("button", { name: "forgotPassword.submit" }),
    );

    await expect(
      screen.findByText("forgotPassword.invalidEmail"),
    ).resolves.toBeVisible();
    expect(emailInput).toHaveAttribute(
      "aria-describedby",
      "forgot-password-email-error",
    );
  });

  it("shows the generic confirmation after a successful request", async () => {
    renderRoute(renderRouteOptions);
    await userEvent.type(
      screen.getByLabelText(/forgotPassword\.email/u),
      "person@example.com",
    );
    await userEvent.click(
      screen.getByRole("button", { name: "forgotPassword.submit" }),
    );

    await expect(
      screen.findByText("forgotPassword.confirmation"),
    ).resolves.toBeVisible();
  });

  it("shows an explicit error when reset requests are rate limited", async () => {
    server.use(postForgotPassword400);
    renderRoute(renderRouteOptions);
    await userEvent.type(
      screen.getByLabelText(/forgotPassword\.email/u),
      "person@example.com",
    );
    await userEvent.click(
      screen.getByRole("button", { name: "forgotPassword.submit" }),
    );

    await expect(
      screen.findByText("forgotPassword.error.PASSWORD_RESET_TOO_FREQUENT"),
    ).resolves.toBeVisible();
  });
});
