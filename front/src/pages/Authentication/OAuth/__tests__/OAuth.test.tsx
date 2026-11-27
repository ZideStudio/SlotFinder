import { getGetV1AuthProviderUrlUrl } from "@Front/api/generated/authentication/authentication";
import { AuthenticationContextProvider } from "@Front/contexts/AuthenticationContext/AuthenticationContextProvider";
// oxlint-disable-next-line import/no-namespace
import * as useAuthenticationContext from "@Front/hooks/useAuthenticationContext";
import { LoaderProvider } from "@Front/providers/loaderProvider/LoaderProvider";
import { renderWithQueryClient } from "@Front/utils/testsUtils/customRender/customRender";
import { screen } from "@testing-library/react";
import { oauthProvidersData } from "../constants";
import { OAuth } from "../OAuth";

const renderOAuth = () =>
  renderWithQueryClient(
    <LoaderProvider>
      <AuthenticationContextProvider>
        <OAuth />
      </AuthenticationContextProvider>
    </LoaderProvider>,
  );

describe("OAuth", () => {
  afterEach(() => {
    vi.restoreAllMocks();
  });

  it("should render all OAuth providers as links with correct aria-labels and generated URLs", () => {
    renderOAuth();

    for (const provider of oauthProvidersData) {
      const link = screen.getByRole("link", {
        name: `authentication.${provider.ariaLabel}`,
      });
      expect(link).toBeInTheDocument();
      expect(link).toHaveAttribute(
        "href",
        `${import.meta.env.FRONT_BACKEND_URL}${getGetV1AuthProviderUrlUrl(provider.id, { returnUrl: "/" })}`,
      );
      expect(link).toHaveAttribute("rel", "noopener noreferrer");
    }
  });

  it("should render all OAuth providers as links with correct generated URLs from postAuthRedirectPath", () => {
    vi.spyOn(
      useAuthenticationContext,
      "useAuthenticationContext",
    ).mockReturnValue({
      isAuthenticated: false,
      postAuthRedirectPath: "/custom-path",
      setPostAuthRedirectPath: vi.fn(),
      resetPostAuthRedirectPath: vi.fn(),
      checkAuthentication: vi.fn(),
      authenticationError: undefined,
    });

    renderOAuth();

    for (const provider of oauthProvidersData) {
      const link = screen.getByRole("link", {
        name: `authentication.${provider.ariaLabel}`,
      });
      expect(link).toBeInTheDocument();
      expect(link).toHaveAttribute(
        "href",
        `${import.meta.env.FRONT_BACKEND_URL}${getGetV1AuthProviderUrlUrl(provider.id, { returnUrl: "/custom-path" })}`,
      );
    }
  });

  it("should render provider icons with aria-hidden", () => {
    renderOAuth();

    for (const provider of oauthProvidersData) {
      const link = screen.getByRole("link", {
        name: `authentication.${provider.ariaLabel}`,
      });
      const icon = link.querySelector("svg");
      expect(icon).toHaveAttribute("aria-hidden");
    }
  });
});
