import { appRoutes } from '@Front/routing/appRoutes';
import { renderRoute, type RenderRouteOptions } from '@Front/utils/testsUtils/customRender/customRender';
import { getAuthStatus400 } from '@Mocks/handlers/authStatusHandlers';
import { server } from '@Mocks/server';
import { screen } from '@testing-library/react';
import { describe, it, expect, beforeEach } from 'vitest';
import { routeObject } from '../../../routing/routes';

const renderRouteOptions: RenderRouteOptions = {
  routes: routeObject,
  routesOptions: { initialEntries: [appRoutes.home()] },
};

describe('Layout', () => {
  describe('when the route has no hideHeader handle', () => {
    it('should render the header banner', async () => {
      renderRoute(renderRouteOptions);

      await expect(screen.findByRole('banner')).resolves.toBeInTheDocument();
    });
  });

  describe('when the route has hideHeader: true', () => {
    beforeEach(() => {
      server.use(getAuthStatus400);
    });

    it('should not render the header banner', async () => {
      renderRoute(renderRouteOptions);

      await screen.findByRole('heading', { level: 1 });

      expect(screen.queryByRole('banner')).not.toBeInTheDocument();
    });
  });
});
