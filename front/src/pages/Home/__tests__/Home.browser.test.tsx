import { appRoutes } from "@Front/routing/appRoutes";
import { renderBrowserRoute } from "@Front/utils/testsUtils/customRender/customRender.browser";
import { worker } from "@Mocks/browser";
import {
  getAuthStatus200,
  getAuthStatus401,
} from "@Mocks/handlers/authStatusHandlers";
import { page } from "vitest/browser";

describe("Home page", () => {
  it("displays the Welcome heading when auth status returns 401", async () => {
    worker.use(getAuthStatus401);

    await renderBrowserRoute({ initialEntry: appRoutes.home() });

    await expect
      .element(
        page.getByRole("heading", {
          level: 1,
          name: "Welcome to the Home page!",
        }),
      )
      .toBeVisible();
  });

  it("displays the Dashboard heading when auth status returns 200", async () => {
    worker.use(getAuthStatus200);

    await renderBrowserRoute({ initialEntry: appRoutes.home() });

    await expect
      .element(page.getByRole("heading", { level: 1, name: "Dashboard" }))
      .toBeVisible();
  });
});
