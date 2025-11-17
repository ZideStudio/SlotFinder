import { appRoutes } from '@Front/routing/appRoutes';
import { renderRoute, type RenderRouteOptions } from '@Front/utils/testsUtils/customRender/customRender';
import { getAuthStatus400 } from '@Mocks/handlers/authStatusHandlers';
import { server } from '@Mocks/server';
import { screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { homeRoutes } from '../../routes';

const renderRouteOptions: RenderRouteOptions = {
  routes: [homeRoutes],
  routesOptions: { initialEntries: [appRoutes.home()] },
};

describe('Welcome', () => {
  beforeEach(() => {
    server.use(getAuthStatus400);
    renderRoute(renderRouteOptions);
  });

  it('renders the home heading', async () => {
    expect(await screen.findByRole('heading', { level: 1, name: 'welcome.title' })).toBeInTheDocument();
  });

  it('renders the sign up link', async () => {
    expect(await screen.findByRole('link', { name: 'Sign Up' })).toHaveAttribute('href', appRoutes.signUp());
  });
});
