import { appRoutes } from '@Front/routing/appRoutes';
import { renderRoute, type RenderRouteOptions } from '@Front/utils/testsUtils/customRender/customRender';
import { screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import { homeRoutes } from '../../routes';

const renderRouteOptions: RenderRouteOptions = {
  routes: [homeRoutes],
  routesOptions: { initialEntries: [appRoutes.home()] },
};

describe('Dashboard', () => {
  it('should render the dashboard heading', async () => {
    renderRoute(renderRouteOptions);

    await expect(screen.findByRole('heading', { level: 1, name: 'dashboard.title' })).resolves.toBeInTheDocument();
  });

  it('should render the header banner', async () => {
    renderRoute(renderRouteOptions);

    await expect(screen.findByRole('banner')).resolves.toBeInTheDocument();
  });
});
