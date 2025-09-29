import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { render, type RenderOptions } from '@testing-library/react';
import type { ComponentProps, ReactNode } from 'react';
import { createMemoryRouter, RouterProvider, type MemoryRouterOpts, type RouteObject } from 'react-router-dom';

export type RenderWithQueryClientOptions = {
  renderOptions?: Omit<RenderOptions, 'queries'>;
  queryClientProviderOptions?: Partial<ComponentProps<typeof QueryClientProvider>>;
};

export type RenderRouteOptions = {
  routes: RouteObject[];
  routesOptions?: MemoryRouterOpts;
  renderOptions?: Omit<RenderOptions, 'queries'>;
  queryClientProviderOptions?: Partial<ComponentProps<typeof QueryClientProvider>>;
};

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
  routes,
  routesOptions,
  renderOptions,
  queryClientProviderOptions,
}: RenderRouteOptions) => {
  const router = createMemoryRouter(routes, routesOptions);

  return renderWithQueryClient(
    <RouterProvider router={router} />,
    { queryClientProviderOptions, renderOptions },
  );
};
