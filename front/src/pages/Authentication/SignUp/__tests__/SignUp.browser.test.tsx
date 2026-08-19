import { screen, waitFor } from "@testing-library/dom";
// oxlint-disable-next-line import/no-namespace
import * as authenticationContextHook from "@Front/hooks/useAuthenticationContext";
import { appRoutes } from "@Front/routing/appRoutes";
import {
  renderBrowserRoute,
  type RenderBrowserRouteOptions,
} from "@Front/utils/testsUtils/customRender/customRender.browser";
import { postAccount400Fixture } from "@Mocks/fixtures/accountFixtures";
import {
  postAccount201,
  postAccount400,
} from "@Mocks/handlers/accountHandlers";
import { worker } from "@Mocks/browser";
import { userEvent } from "@vitest/browser/context";
import { authenticationRoutes } from "../../routes";

const renderRouteOptions: RenderBrowserRouteOptions = {
  routes: [authenticationRoutes],
  routesOptions: { initialEntries: [appRoutes.signUp()] },
};

describe("SignUp", () => {
  beforeEach(() => {
    worker.use(postAccount201);
  });

  afterEach(() => {
    worker.resetHandlers();
  });

  it("renders all form fields, submit button and oauth", async () => {
    await renderBrowserRoute(renderRouteOptions);

    expect(screen.getByLabelText("signUp.username")).toBeInTheDocument();
    expect(screen.getByLabelText("signUp.email")).toBeInTheDocument();
    expect(screen.getByLabelText("signUp.password")).toBeInTheDocument();
    expect(screen.getByLabelText("signUp.confirmPassword")).toBeInTheDocument();
    expect(
      screen.getByRole("button", { name: "signUp.submit" }),
    ).toBeInTheDocument();
    expect(
      screen.getByRole("navigation", {
        name: "authentication.signInWithProvider",
      }),
    ).toBeInTheDocument();
  });

  it("shows validation errors for empty fields", async () => {
    await renderBrowserRoute(renderRouteOptions);

    await userEvent.click(
      screen.getByRole("button", { name: "signUp.submit" }),
    );

    await expect(
      screen.findByText("signUp.requiredUsername"),
    ).resolves.toBeInTheDocument();
    expect(screen.getByText("signUp.requiredEmail")).toBeInTheDocument();
    expect(screen.getByText("signUp.requiredPassword")).toBeInTheDocument();
    expect(
      screen.getByText("signUp.requiredConfirmPassword"),
    ).toBeInTheDocument();
  });

  it.each([
    {
      password: "1234567",
      expectedError: 'signUp.minLengthPassword::{"min":8}',
      description: "minimum length error",
    },
    {
      password: "12345678!",
      expectedError: "signUp.passwordComplexity",
      description: "must contain letters error",
    },
    {
      password: "password1!",
      expectedError: "signUp.passwordComplexity",
      description: "must contain numbers error",
    },
    {
      password: "Password!",
      expectedError: "signUp.passwordComplexity",
      description: "must contain numbers error",
    },
    {
      password: "Password1",
      expectedError: "signUp.passwordComplexity",
      description: "must contain symbols error",
    },
  ])("shows password $description", async ({ password, expectedError }) => {
    await renderBrowserRoute(renderRouteOptions);

    await userEvent.type(screen.getByLabelText("signUp.username"), "testuser");
    await userEvent.type(
      screen.getByLabelText("signUp.email"),
      "test@example.com",
    );
    await userEvent.type(screen.getByLabelText("signUp.password"), password);
    await userEvent.type(
      screen.getByLabelText("signUp.confirmPassword"),
      password,
    );
    await userEvent.click(
      screen.getByRole("button", { name: "signUp.submit" }),
    );

    await expect(screen.findByText(expectedError)).resolves.toBeInTheDocument();
  });

  it("shows accessible error when confirm password field does not match with password field", async () => {
    await renderBrowserRoute(renderRouteOptions);

    await userEvent.type(screen.getByLabelText("signUp.username"), "testuser");
    await userEvent.type(
      screen.getByLabelText("signUp.email"),
      "test@example.com",
    );
    await userEvent.type(
      screen.getByLabelText("signUp.password"),
      "Password1!",
    );
    await userEvent.type(
      screen.getByLabelText("signUp.confirmPassword"),
      "DifferentPassword1!",
    );
    await userEvent.click(
      screen.getByRole("button", { name: "signUp.submit" }),
    );

    const confirmPasswordError = await screen.findByRole("alert");
    expect(confirmPasswordError).toHaveTextContent(
      "signUp.passwordsDoNotMatch",
    );
    expect(confirmPasswordError).toHaveAttribute("id", "confirmPassword-error");

    const confirmPasswordInput = screen.getByLabelText(
      "signUp.confirmPassword",
    );
    expect(confirmPasswordInput).toHaveAttribute(
      "aria-describedby",
      "confirmPassword-error",
    );
  });

  it("checks authentication from authentication context on successful submission", async () => {
    const checkAuthentication = vi.fn();
    vi.spyOn(
      authenticationContextHook,
      "useAuthenticationContext",
    ).mockReturnValue({
      checkAuthentication,
      isAuthenticated: undefined,
      postAuthRedirectPath: undefined,
      setPostAuthRedirectPath: vi.fn(),
      resetPostAuthRedirectPath: vi.fn(),
    });

    await renderBrowserRoute(renderRouteOptions);

    await userEvent.type(screen.getByLabelText("signUp.username"), "testuser");
    await userEvent.type(
      screen.getByLabelText("signUp.email"),
      "test@example.com",
    );
    await userEvent.type(
      screen.getByLabelText("signUp.password"),
      "Password1!",
    );
    await userEvent.type(
      screen.getByLabelText("signUp.confirmPassword"),
      "Password1!",
    );
    await userEvent.click(
      screen.getByRole("button", { name: "signUp.submit" }),
    );

    await waitFor(() => {
      expect(checkAuthentication).toHaveBeenCalledTimes(1);
    });
  });
});

describe("SignUp error handling", () => {
  beforeEach(() => {
    worker.use(postAccount400);
  });

  it("shows error message on failed submission", async () => {
    await renderBrowserRoute(renderRouteOptions);

    await userEvent.type(screen.getByLabelText("signUp.username"), "failuser");
    await userEvent.type(
      screen.getByLabelText("signUp.email"),
      "fail@example.com",
    );
    await userEvent.type(
      screen.getByLabelText("signUp.password"),
      "Password1!",
    );
    await userEvent.type(
      screen.getByLabelText("signUp.confirmPassword"),
      "Password1!",
    );
    await userEvent.click(
      screen.getByRole("button", { name: "signUp.submit" }),
    );

    await expect(
      screen.findByText(`signUp.error.${postAccount400Fixture.code}`),
    ).resolves.toBeInTheDocument();
  });
});
