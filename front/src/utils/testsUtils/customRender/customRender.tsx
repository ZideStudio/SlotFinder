import { AuthenticationContextProvider } from '@Front/contexts/AuthenticationContext/AuthenticationContextProvider';
import { routeObject } from '@Front/routing/routes';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { render, type RenderOptions } from '@testing-library/react';
import type { ComponentProps, ReactNode } from 'react';
import { createMemoryRouter, RouterProvider, type MemoryRouterOpts, type RouteObject } from 'react-router';

export type RenderWithQueryClientOptions = {
  renderOptions?: Omit<RenderOptions, 'queries'>;
  queryClientProviderOptions?: Partial<ComponentProps<typeof QueryClientProvider>>;
};

export type RenderRouteOptions = {
  routes?: RouteObject[];
  renderOptions?: Omit<RenderOptions, 'queries'>;
  queryClientProviderOptions?: Partial<ComponentProps<typeof QueryClientProvider>>;
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
  { queryClientProviderOptions, renderOptions }: RenderWithQueryClientOptions = {},
) => {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
        refetchOnWindowFocus: false,
        refetchOnReconnect: false,
      },
    },
  });

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
  const routesOptionsWithEntry = {
    initialEntries: initialEntry ? [initialEntry] : undefined,
    ...routesOptions,
  };
  const router = createMemoryRouter(routes, routesOptionsWithEntry);

  return renderWithQueryClient(
    <AuthenticationContextProvider>
      <RouterProvider router={router} />
    </AuthenticationContextProvider>,
    { queryClientProviderOptions, renderOptions },
  );
};
