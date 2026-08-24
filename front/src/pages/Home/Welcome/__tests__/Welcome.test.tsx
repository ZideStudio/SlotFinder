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
      screen.findByRole("heading", { level: 1, name: "welcome.slotFinder" }),
    ).resolves.toBeInTheDocument();
  });

  it("does not render the header banner", () => {
    expect(screen.queryByRole("banner")).toBeNull();
  });

  it("does render the connexion button", async () => {
    await expect(
      screen.findByRole("button", { name: "welcome.connexion" }),
    ).resolves.toBeInTheDocument();
  });

  it("does render the create an event button", async () => {
    await expect(
      screen.findByRole("button", { name: "welcome.create an event" }),
    ).resolves.toBeInTheDocument();
  });
});
