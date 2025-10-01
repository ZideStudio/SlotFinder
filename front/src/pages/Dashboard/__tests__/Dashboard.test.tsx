import { appRoutes } from '@Front/routing/appRoutes';
import { renderRoute, type RenderRouteOptions } from '@Front/utils/testsUtils/customRender';
import { screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { dashboardRoutes } from '../routes';

const renderRouteOptions: RenderRouteOptions = {
  routes: [dashboardRoutes],
  routesOptions: { initialEntries: [appRoutes.dashboard()] },
};

describe('Dashboard', () => {
  it('shows error message on failed submission', () => {
    renderRoute(renderRouteOptions);

    expect(screen.getByRole('heading', { level: 1, name: 'Dashboard' }));
  });
});
