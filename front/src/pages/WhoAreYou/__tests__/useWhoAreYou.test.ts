// oxlint-disable-next-line import/no-namespace
import * as useAuthenticationContext from "@Front/hooks/useAuthenticationContext";
import { TestProviders } from "@Front/utils/testsUtils/customRender/TestProviders";
import {
  getAccountMe200,
  patchAccount200,
  patchAccount400,
  patchAvatarAccount200,
  patchAvatarAccount400,
} from "@Mocks/handlers/accountHandlers";
import { server } from "@Mocks/server";
import { renderHook, waitFor } from "@testing-library/react";
import { useWhoAreYou } from "../useWhoAreYou";

const createAvatarFileList = () => {
  const avatar = new File(["avatar"], "avatar.png", { type: "image/png" });
  return [avatar] as unknown as FileList;
};

const renderHookWithProviders = (hook: () => ReturnType<typeof useWhoAreYou>) =>
  renderHook(hook, { wrapper: TestProviders });

describe("useWhoAreYou - success scenarios", () => {
  const originalTemporal = globalThis.Temporal;

  beforeAll(() => {
    globalThis.Temporal = {
      Now: {
        timeZoneId: () => "Europe/Paris",
      },
    } as unknown as typeof Temporal;
  });

  beforeEach(() => {
    server.use(getAccountMe200, patchAccount200, patchAvatarAccount200);
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  afterAll(() => {
    globalThis.Temporal = originalTemporal;
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
    });
  });

  it("should refresh authentication state after a successful submit", async () => {
    const setError = vi.fn();
    const reset = vi.fn();
    const checkAuthentication = vi.fn();

    vi.spyOn(
      useAuthenticationContext,
      "useAuthenticationContext",
    ).mockReturnValue({
      checkAuthentication,
    } as unknown as ReturnType<
      typeof useAuthenticationContext.useAuthenticationContext
    >);

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
      expect(checkAuthentication).toHaveBeenCalledTimes(1);
    });
  });
});

describe("useWhoAreYou - error scenarios", () => {
  const originalTemporal = globalThis.Temporal;

  beforeAll(() => {
    globalThis.Temporal = {
      Now: {
        timeZoneId: () => "Europe/Paris",
      },
    } as unknown as typeof Temporal;
  });

  beforeEach(() => {
    server.resetHandlers();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  afterAll(() => {
    globalThis.Temporal = originalTemporal;
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
