import { appRoutes } from "@Front/routing/appRoutes";
import {
  renderRoute,
  type RenderRouteOptions,
} from "@Front/utils/testsUtils/customRender/customRender";
import { getAuthStatus401 } from "@Mocks/handlers/authStatusHandlers";
import { server } from "@Mocks/server";
import { screen } from "@testing-library/react";
import { homeRoutes } from "../../routes";

const renderRouteOptions: RenderRouteOptions = {
  routes: [homeRoutes],
  routesOptions: { initialEntries: [appRoutes.home()] },
};

describe("Welcome", () => {
  beforeEach(() => {
    server.use(getAuthStatus401);
    renderRoute(renderRouteOptions);
  });

  it("renders the home heading", async () => {
    await expect(
      screen.findByRole("heading", { level: 1, name: "welcome.title" }),
    ).resolves.toBeInTheDocument();
  });

  it("renders the sign up link", async () => {
    await expect(
      screen.findByRole("link", { name: "Sign Up" }),
    ).resolves.toHaveAttribute("href", appRoutes.signUp());
  });

  it("does not render the header banner", () => {
    expect(screen.queryByRole("banner")).toBeNull();
  });
});
