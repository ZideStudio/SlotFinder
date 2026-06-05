import { routeObject } from "@Front/routing/routes";
import {
  createMemoryRouter,
  RouterProvider,
  type MemoryRouterOpts,
  type RouteObject,
} from "react-router";
import { render } from "vitest-browser-react";
import { createQueryClient, TestProviders } from "./TestProviders";

export type RenderBrowserRouteOptions = {
  routes?: RouteObject[];
} & (
  | {
      initialEntry: string;
      routesOptions?: MemoryRouterOpts;
    }
  | {
      initialEntry?: string;
      routesOptions: MemoryRouterOpts;
    }
);

export const renderBrowserRoute = ({
  initialEntry,
  routes = routeObject,
  routesOptions = {},
}: RenderBrowserRouteOptions) => {
  const router = createMemoryRouter(routes, {
    initialEntries: initialEntry ? [initialEntry] : undefined,
    ...routesOptions,
  });

  return render(
    <TestProviders client={createQueryClient()}>
      <RouterProvider router={router} />
    </TestProviders>,
  );
};
