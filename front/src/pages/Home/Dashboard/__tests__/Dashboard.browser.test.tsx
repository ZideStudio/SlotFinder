import { appRoutes } from "@Front/routing/appRoutes";
import { renderBrowserRoute } from "@Front/utils/testsUtils/customRender/customRender.browser";
import { worker } from "@Mocks/browser";
import { getAuthStatus200 } from "@Mocks/handlers/authStatusHandlers";
import { page } from "vitest/browser";

describe("Dashboard Page", () => {
  beforeEach(() => {
    worker.use(getAuthStatus200);
  });

  it("should render the dashboard heading", async () => {
    await renderBrowserRoute({ initialEntry: appRoutes.home() });

    await expect
      .element(page.getByRole("heading", { level: 1, name: "My events" }))
      .toBeInTheDocument();
  });

  it("does render the header banner", async () => {
    await renderBrowserRoute({ initialEntry: appRoutes.home() });

    await expect.element(page.getByRole("banner")).toBeInTheDocument();
  });

  it("does render no events message when there are no events", async () => {
    await renderBrowserRoute({ initialEntry: appRoutes.home() });

    await expect.element(page.getByText("No events here")).toBeInTheDocument();
  });
});
