import { AuthenticationContextProvider } from "@Front/contexts/AuthenticationContext/AuthenticationContextProvider";
import { routeObject } from "@Front/routing/routes";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import {
  createMemoryRouter,
  RouterProvider,
  type MemoryRouterOpts,
  type RouteObject,
} from "react-router";
import { render } from "vitest-browser-react";

const createQueryClient = () =>
  new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
        refetchOnWindowFocus: false,
        refetchOnReconnect: false,
      },
    },
  });

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
    <QueryClientProvider client={createQueryClient()}>
      <AuthenticationContextProvider>
        <RouterProvider router={router} />
      </AuthenticationContextProvider>
    </QueryClientProvider>,
  );
};
