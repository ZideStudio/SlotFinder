import { appRoutes } from "@Front/routing/appRoutes";
import { screen } from "@testing-library/react";
import { renderBrowserRoute } from "@Front/utils/testsUtils/customRender/customRender.browser";

describe("Dashboard Page", () => {
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
});
