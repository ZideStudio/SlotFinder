import { routeObject } from "@Front/routing/routes";
import { QueryClientProvider } from "@tanstack/react-query";
import { render, type RenderOptions } from "@testing-library/react";
import type { ComponentProps, ReactNode } from "react";
import {
  createMemoryRouter,
  RouterProvider,
  type MemoryRouterOpts,
  type RouteObject,
} from "react-router";
import { createQueryClient, TestProviders } from "./TestProviders";

export type RenderWithQueryClientOptions = {
  renderOptions?: Omit<RenderOptions, "queries">;
  queryClientProviderOptions?: Partial<
    ComponentProps<typeof QueryClientProvider>
  >;
};

export type RenderRouteOptions = {
  routes?: RouteObject[];
  renderOptions?: Omit<RenderOptions, "queries">;
  queryClientProviderOptions?: Partial<
    ComponentProps<typeof QueryClientProvider>
  >;
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

export const renderWithQueryClient = (
  ui: ReactNode,
  {
    queryClientProviderOptions,
    renderOptions,
  }: RenderWithQueryClientOptions = {},
) => {
  const queryClient = createQueryClient();

  return render(
    <QueryClientProvider client={queryClient} {...queryClientProviderOptions}>
      {ui}
    </QueryClientProvider>,
    renderOptions,
  );
};

export const renderRoute = ({
  initialEntry,
  routes = routeObject,
  routesOptions = {},
  renderOptions,
  queryClientProviderOptions,
}: RenderRouteOptions) => {
  const router = createMemoryRouter(routes, {
    initialEntries: initialEntry ? [initialEntry] : undefined,
    ...routesOptions,
  });

  const client = queryClientProviderOptions?.client ?? createQueryClient();

  return render(
    <TestProviders client={client}>
      <RouterProvider router={router} />
    </TestProviders>,
    renderOptions,
  );
};
