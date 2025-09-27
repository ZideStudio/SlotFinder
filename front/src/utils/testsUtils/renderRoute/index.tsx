import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { render, type RenderOptions } from '@testing-library/react';
import type { ComponentProps } from 'react';
import { createMemoryRouter, RouterProvider, type MemoryRouterOpts, type RouteObject } from 'react-router';

export type RenderRouteOptions = {
  routes: RouteObject[];
  routesOptions?: MemoryRouterOpts;
  renderOptions?: Omit<RenderOptions, 'queries'>;
  queryClientProviderOptions?: Partial<ComponentProps<typeof QueryClientProvider>>;
};

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: false,
      refetchOnWindowFocus: false,
      refetchOnReconnect: false,
    },
  },
});

export const renderRoute = ({
  routes,
  routesOptions,
  renderOptions,
  queryClientProviderOptions,
}: RenderRouteOptions) => {
  const router = createMemoryRouter(routes, routesOptions);

  return render(
    <QueryClientProvider client={queryClient} {...queryClientProviderOptions}>
      <RouterProvider router={router} />
    </QueryClientProvider>,
    renderOptions,
  );
};
