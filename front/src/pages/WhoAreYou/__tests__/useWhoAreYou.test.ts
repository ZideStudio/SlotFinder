import { TestProviders } from "@Front/utils/testsUtils/customRender/TestProviders";
import {
  patchAccount200,
  patchAccount400,
  patchAvatarAccount200,
  patchAvatarAccount400,
} from "@Mocks/handlers/accountHandlers";
import { server } from "@Mocks/server";
import { renderHook, waitFor } from "@testing-library/react";
import { useWhoAreYou } from "../useWhoAreYou";

describe("useWhoAreYou", () => {
  const originalTemporal = globalThis.Temporal;

  beforeAll(() => {
    globalThis.Temporal = {
      Now: {
        timeZoneId: () => "Europe/Paris",
      },
    } as unknown as typeof Temporal;
  });

  beforeEach(() => {
    server.use(patchAccount200, patchAvatarAccount200);
  });

  afterAll(() => {
    globalThis.Temporal = originalTemporal;
  });

  const createAvatarFileList = () => {
    const avatar = new File(["avatar"], "avatar.png", { type: "image/png" });
    return [avatar] as unknown as FileList;
  };

  const renderHookWithProviders = (
    hook: () => ReturnType<typeof useWhoAreYou>,
  ) => renderHook(hook, { wrapper: TestProviders });

  it("should submit account and avatar payloads on successful submit", async () => {
    const setError = vi.fn();

    const { result } = renderHookWithProviders(() =>
      useWhoAreYou({ setError }),
    );

    await result.current.handleSubmit({
      avatar: createAvatarFileList(),
      username: "john_doe",
      color: "#ff0000",
      termsAccepted: true,
    });

    expect(setError).not.toHaveBeenCalled();
    expect(result.current.submitError).toBeNull();
  });

  it("should handle username already taken error from API", async () => {
    server.use(patchAccount400("USERNAME_ALREADY_TAKEN"));
    const setError = vi.fn();

    const { result } = renderHookWithProviders(() =>
      useWhoAreYou({ setError }),
    );

    await result.current.handleSubmit({
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
    server.use(patchAccount400("INVALID_COLOR_FORMAT"));
    const setError = vi.fn();

    const { result } = renderHookWithProviders(() =>
      useWhoAreYou({ setError }),
    );

    await result.current.handleSubmit({
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
    server.use(patchAccount200, patchAvatarAccount400);
    const setError = vi.fn();

    const { result } = renderHookWithProviders(() =>
      useWhoAreYou({ setError }),
    );

    await result.current.handleSubmit({
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

  it("should handle account patch error and allow retry", async () => {
    server.use(patchAccount400("INVALID_COLOR_FORMAT"));
    const setError = vi.fn();

    const { result } = renderHookWithProviders(() =>
      useWhoAreYou({ setError }),
    );

    // First submit with error
    await result.current.handleSubmit({
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

    // Reset for successful retry
    server.use(patchAccount200, patchAvatarAccount200);
    vi.clearAllMocks();

    // Retry successfully
    await result.current.handleSubmit({
      avatar: createAvatarFileList(),
      username: "john_doe",
      color: "#ff0000",
      termsAccepted: true,
    });

    expect(result.current.submitError).toBeNull();
    expect(setError).not.toHaveBeenCalled();
  });
});
