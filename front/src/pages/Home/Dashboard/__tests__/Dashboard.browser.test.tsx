import { appRoutes } from "@Front/routing/appRoutes";
import { screen } from "@testing-library/react";
import { renderBrowserRoute } from "@Front/utils/testsUtils/customRender/customRender.browser";
import { worker } from "@Mocks/browser";
import { getAuthStatus200 } from "@Mocks/handlers/authStatusHandlers";

describe("Dashboard Page", () => {
  beforeEach(() => {
    worker.use(getAuthStatus200);
  });
  it("should render the dashboard heading", async () => {
    await renderBrowserRoute({ initialEntry: appRoutes.home() });

    await expect(
      screen.findByRole("heading", {
        level: 1,
        name: "My events",
      }),
    ).resolves.toBeInTheDocument();
  });

  it("does render the header banner", async () => {
    await renderBrowserRoute({ initialEntry: appRoutes.home() });

    await expect(screen.findByRole("banner")).resolves.toBeInTheDocument();
  });

  it("does render no events message when there are no events", async () => {
    await renderBrowserRoute({ initialEntry: appRoutes.home() });

    await expect(
      screen.findByText("No events here"),
    ).resolves.toBeInTheDocument();
  });
});
