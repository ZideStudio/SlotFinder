import { AuthenticationContext } from "@Front/contexts/AuthenticationContext/AuthenticationContext";
import type { AuthenticationContextType } from "@Front/contexts/AuthenticationContext/types";
import { createQueryClient } from "@Front/utils/testsUtils/customRender/TestProviders";
// oxlint-disable-next-line import/no-namespace
import * as patchAccountHook from "@Front/api/account/patchAccount/usePatchAccount";
// oxlint-disable-next-line import/no-namespace
import * as patchAccountAvatarHook from "@Front/api/account/patchAccountAvatar/usePatchAccountAvatar";
import {
  getAccountMe200,
  patchAccount200,
  patchAccount400,
  patchAvatarAccount200,
  patchAvatarAccount400,
} from "@Mocks/handlers/accountHandlers";
import { server } from "@Mocks/server";
import { QueryClientProvider } from "@tanstack/react-query";
import { act, renderHook, waitFor } from "@testing-library/react";
import { createElement, type PropsWithChildren } from "react";
import { useWhoAreYou } from "../useWhoAreYou";

const createAvatarFileList = () => {
  const avatar = new File(["avatar"], "avatar.png", { type: "image/png" });
  return [avatar] as unknown as FileList;
};

const createAuthenticationContextValue = (
  overrides: Partial<AuthenticationContextType> = {},
): AuthenticationContextType => ({
  isAuthenticated: true,
  authenticationError: undefined,
  checkAuthentication: vi.fn(),
  postAuthRedirectPath: undefined,
  setPostAuthRedirectPath: vi.fn(),
  resetPostAuthRedirectPath: vi.fn(),
  ...overrides,
});

const renderHookWithProviders = (
  hook: () => ReturnType<typeof useWhoAreYou>,
  authContextOverrides: Partial<AuthenticationContextType> = {},
) => {
  const client = createQueryClient();
  const authContextValue = createAuthenticationContextValue(authContextOverrides);

  const wrapper = ({ children }: PropsWithChildren) =>
    createElement(
      QueryClientProvider,
      { client },
      createElement(
        AuthenticationContext.Provider,
        { value: authContextValue },
        children,
      ),
    );

  return {
    ...renderHook(hook, { wrapper }),
    authContextValue,
  };
};

describe("useWhoAreYou - success scenarios", () => {
  beforeEach(() => {
    server.use(getAccountMe200, patchAccount200, patchAvatarAccount200);
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  afterAll(() => {
    server.resetHandlers();
  });

  it("should submit account and avatar payloads on successful submit", () => {
    const setError = vi.fn();
    const reset = vi.fn();

    const { result } = renderHookWithProviders(() =>
      useWhoAreYou({ setError, reset }),
    );

    result.current.handleSubmit({
      avatar: createAvatarFileList(),
      username: "john_doe",
      color: "#ff0000",
      termsAccepted: true,
    });

    expect(setError).not.toHaveBeenCalled();
    expect(result.current.submitError).toBeNull();
  });

  it("should return defaultAvatarUrl from account data", async () => {
    const setError = vi.fn();
    const reset = vi.fn();

    const { result } = renderHookWithProviders(() =>
      useWhoAreYou({ setError, reset }),
    );

    await waitFor(() => {
      expect(result.current.defaultAvatarUrl).toBe(
        "/api/v1/account/123456789/avatar",
      );
      expect(result.current.hasAcceptedCurrentTermsVersion).toBe(true);
    });
  });

  it("should prefill termsAccepted as true when current terms are already accepted", async () => {
    const setError = vi.fn();
    const reset = vi.fn();

    renderHookWithProviders(() => useWhoAreYou({ setError, reset }));

    await waitFor(() => {
      expect(reset).toHaveBeenCalledWith(
        expect.objectContaining({ termsAccepted: true }),
      );
    });
  });

  it("should refresh authentication state after a successful submit", async () => {
    const setError = vi.fn();
    const reset = vi.fn();
    const checkAuthentication = vi.fn();
    const patchAccount = vi.fn().mockResolvedValue({});
    const patchAccountAvatar = vi.fn().mockResolvedValue(undefined);

    vi.spyOn(patchAccountHook, "usePatchAccount").mockReturnValue({
      patchAccount,
      isLoading: false,
      errorCode: undefined,
    });
    vi.spyOn(
      patchAccountAvatarHook,
      "usePatchAccountAvatar",
    ).mockReturnValue({
      patchAccountAvatar,
      isLoading: false,
    });

    const { result } = renderHookWithProviders(() =>
      useWhoAreYou({ setError, reset }),
      { checkAuthentication },
    );

    await act(async () => {
      await result.current.handleSubmit({
        avatar: createAvatarFileList(),
        username: "john_doe",
        color: "#ff0000",
        termsAccepted: true,
      });
    });

    await waitFor(() => {
      expect(patchAccount).toHaveBeenCalledTimes(1);
      expect(patchAccountAvatar).toHaveBeenCalledTimes(1);
      expect(checkAuthentication).toHaveBeenCalledTimes(1);
    });
  });
});

describe("useWhoAreYou - error scenarios", () => {
  beforeEach(() => {
    server.resetHandlers();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  afterAll(() => {
    server.resetHandlers();
  });

  it("should handle username already taken error from API", async () => {
    server.use(
      getAccountMe200,
      patchAccount400("USERNAME_ALREADY_TAKEN"),
      patchAvatarAccount200,
    );

    const setError = vi.fn();
    const reset = vi.fn();

    const { result } = renderHookWithProviders(() =>
      useWhoAreYou({ setError, reset }),
    );

    result.current.handleSubmit({
      avatar: createAvatarFileList(),
      username: "taken_username",
      color: "#ff0000",
      termsAccepted: true,
    });

    await waitFor(() => {
      expect(setError).toHaveBeenCalledWith("username", {
        message: "whoAreYou.error.USERNAME_ALREADY_TAKEN",
      });
    });
  });

  it("should handle invalid color format error from API", async () => {
    server.use(
      getAccountMe200,
      patchAccount400("INVALID_COLOR_FORMAT"),
      patchAvatarAccount200,
    );

    const setError = vi.fn();
    const reset = vi.fn();

    const { result } = renderHookWithProviders(() =>
      useWhoAreYou({ setError, reset }),
    );

    result.current.handleSubmit({
      avatar: createAvatarFileList(),
      username: "john_doe",
      color: "invalid",
      termsAccepted: true,
    });

    await waitFor(() => {
      expect(setError).toHaveBeenCalledWith("color", {
        message: "whoAreYou.error.INVALID_COLOR_FORMAT",
      });
    });
  });

  it("should handle avatar upload failure from API", async () => {
    server.use(getAccountMe200, patchAccount200, patchAvatarAccount400);

    const setError = vi.fn();
    const reset = vi.fn();

    const { result } = renderHookWithProviders(() =>
      useWhoAreYou({ setError, reset }),
    );

    result.current.handleSubmit({
      avatar: createAvatarFileList(),
      username: "john_doe",
      color: "#ff0000",
      termsAccepted: true,
    });

    await waitFor(() => {
      expect(setError).toHaveBeenCalledWith("avatar", {
        message: "whoAreYou.error.AVATAR_UPLOAD_FAILED",
      });
    });
  });
});
