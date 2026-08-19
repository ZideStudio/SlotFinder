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
import { page } from "vitest/browser";
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

    await expect
      .element(page.getByLabelText("signUp.username"))
      .toBeInTheDocument();
    await expect
      .element(page.getByLabelText("signUp.email"))
      .toBeInTheDocument();
    await expect
      .element(page.getByLabelText("signUp.password"))
      .toBeInTheDocument();
    await expect
      .element(page.getByLabelText("signUp.confirmPassword"))
      .toBeInTheDocument();
    await expect
      .element(page.getByRole("button", { name: "signUp.submit" }))
      .toBeInTheDocument();
    await expect
      .element(
        page.getByRole("navigation", {
          name: "authentication.signInWithProvider",
        }),
      )
      .toBeInTheDocument();
  });

  it("shows validation errors for empty fields", async () => {
    await renderBrowserRoute(renderRouteOptions);

    await userEvent.click(page.getByRole("button", { name: "signUp.submit" }));

    await expect
      .element(page.getByText("signUp.requiredUsername"))
      .toBeInTheDocument();
    await expect
      .element(page.getByText("signUp.requiredEmail"))
      .toBeInTheDocument();
    await expect
      .element(page.getByText("signUp.requiredPassword"))
      .toBeInTheDocument();
    await expect
      .element(page.getByText("signUp.requiredConfirmPassword"))
      .toBeInTheDocument();
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

    await userEvent.type(page.getByLabelText("signUp.username"), "testuser");
    await userEvent.type(
      page.getByLabelText("signUp.email"),
      "test@example.com",
    );
    await userEvent.type(page.getByLabelText("signUp.password"), password);
    await userEvent.type(
      page.getByLabelText("signUp.confirmPassword"),
      password,
    );
    await userEvent.click(page.getByRole("button", { name: "signUp.submit" }));

    await expect.element(page.getByText(expectedError)).toBeInTheDocument();
  });

  it("shows accessible error when confirm password field does not match with password field", async () => {
    await renderBrowserRoute(renderRouteOptions);

    await userEvent.type(page.getByLabelText("signUp.username"), "testuser");
    await userEvent.type(
      page.getByLabelText("signUp.email"),
      "test@example.com",
    );
    await userEvent.type(page.getByLabelText("signUp.password"), "Password1!");
    await userEvent.type(
      page.getByLabelText("signUp.confirmPassword"),
      "DifferentPassword1!",
    );
    await userEvent.click(page.getByRole("button", { name: "signUp.submit" }));

    const confirmPasswordError = page.getByRole("alert");
    await expect
      .element(confirmPasswordError)
      .toHaveTextContent("signUp.passwordsDoNotMatch");
    await expect
      .element(confirmPasswordError)
      .toHaveAttribute("id", "confirmPassword-error");

    await expect
      .element(page.getByLabelText("signUp.confirmPassword"))
      .toHaveAttribute("aria-describedby", "confirmPassword-error");
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

    await userEvent.type(page.getByLabelText("signUp.username"), "testuser");
    await userEvent.type(
      page.getByLabelText("signUp.email"),
      "test@example.com",
    );
    await userEvent.type(page.getByLabelText("signUp.password"), "Password1!");
    await userEvent.type(
      page.getByLabelText("signUp.confirmPassword"),
      "Password1!",
    );
    await userEvent.click(page.getByRole("button", { name: "signUp.submit" }));

    await vi.waitFor(() => {
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

    await userEvent.type(page.getByLabelText("signUp.username"), "failuser");
    await userEvent.type(
      page.getByLabelText("signUp.email"),
      "fail@example.com",
    );
    await userEvent.type(page.getByLabelText("signUp.password"), "Password1!");
    await userEvent.type(
      page.getByLabelText("signUp.confirmPassword"),
      "Password1!",
    );
    await userEvent.click(page.getByRole("button", { name: "signUp.submit" }));

    await expect
      .element(page.getByText(`signUp.error.${postAccount400Fixture.code}`))
      .toBeInTheDocument();
  });
});
