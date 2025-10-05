import { appRoutes } from '@Front/routing/appRoutes';
import { routeObject } from '@Front/routing/routes';
import { renderRoute, type RenderRouteOptions } from '@Front/utils/testsUtils/customRender';
import { getAuthStatus200, getAuthStatus400 } from '@Mocks/handlers/authStatusHandlers';
import { server } from '@Mocks/server';
import { screen } from '@testing-library/react';

const getRenderRouteOptions = (initialEntry: string): RenderRouteOptions => ({
  routes: routeObject,
  routesOptions: { initialEntries: [initialEntry] },
});

afterEach(() => {
  server.resetHandlers();
});

describe('AuthenticationProtection with valid credentials', () => {
  beforeEach(() => {
    server.use(getAuthStatus200);
  });

  it('should render children when route does not require authentication and the user is authenticated', async () => {
    renderRoute(getRenderRouteOptions(appRoutes.home()));

    expect(await screen.findByText('home.welcome')).toBeInTheDocument();
  });

  it('should render children when route requires authentication and the user is authenticated', async () => {
    renderRoute(getRenderRouteOptions(appRoutes.dashboard()));

    expect(await screen.findByText('dashboard.title')).toBeInTheDocument();
  });

  it('should redirect to dashboard when route requires no authentication and the user is authenticated', async () => {
    renderRoute(getRenderRouteOptions(appRoutes.signUp()));

    expect(await screen.findByText('dashboard.title')).toBeInTheDocument();
  });
});

describe('AuthenticationProtection with invalid credentials', () => {
  beforeEach(() => {
    server.use(getAuthStatus400);
  });

  it('should render children when route does not require authentication and the user is not authenticated', async () => {
    renderRoute(getRenderRouteOptions(appRoutes.home()));

    expect(await screen.findByText('home.welcome')).toBeInTheDocument();
  });

  it('should redirect to signUp when route requires authentication and the user is not authenticated', async () => {
    renderRoute(getRenderRouteOptions(appRoutes.dashboard()));

    expect(await screen.findByText('signUp.title')).toBeInTheDocument();
  });

  it('should render children when route requires no authentication and the user is not authenticated', async () => {
    renderRoute(getRenderRouteOptions(appRoutes.signUp()));

    expect(await screen.findByText('signUp.title')).toBeInTheDocument();
  });
});
