import { appRoutes } from '@Front/routing/appRoutes';
import { renderRoute, type RenderRouteOptions } from '@Front/utils/testsUtils/customRender/customRender';
import { screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { dashboardRoutes } from '../routes';

const renderRouteOptions: RenderRouteOptions = {
  routes: [dashboardRoutes],
  routesOptions: { initialEntries: [appRoutes.dashboard()] },
};

describe('Dashboard', () => {
  it('renders the dashboard heading', () => {
    renderRoute(renderRouteOptions);

    expect(screen.getByRole('heading', { level: 1, name: 'dashboard.title' })).toBeInTheDocument();
  });
});
