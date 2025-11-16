import { appRoutes } from '@Front/routing/appRoutes';
import { renderRoute, type RenderRouteOptions } from '@Front/utils/testsUtils/customRender/customRender';
import { screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { homeRoutes } from '../../routes';

const renderRouteOptions: RenderRouteOptions = {
  routes: [homeRoutes],
  routesOptions: { initialEntries: [appRoutes.home()] },
};

describe('Dashboard', () => {
  it('renders the dashboard heading', async () => {
    renderRoute(renderRouteOptions);

    expect(await screen.findByRole('heading', { level: 1, name: 'dashboard.title' })).toBeInTheDocument();
  });
});
