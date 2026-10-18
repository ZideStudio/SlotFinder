import { getRuntimeEnv } from "@Front/config/runtimeEnv";

const enableMock = async () => {
  if (getRuntimeEnv().FRONT_MOCK_ENABLE === "true") {
    const { worker } = await import("@Mocks/browser");

    return worker.start({
      onUnhandledFrame: "warn",
    });
  }
};

// oxlint-disable-next-line jest/require-hook
enableMock().then(() => {
  import("./bootstrap").then(({ default: AppReact }) => {
    customElements.define("app-react", AppReact);
  });
});
