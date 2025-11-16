import * as useAuthenticationContext from '@Front/hooks/useAuthenticationContext';
import { routeObject } from '@Front/routing/routes';
import { renderRoute } from '@Front/utils/testsUtils/customRender/customRender';
import { getAuthStatus200, getAuthStatus400 } from '@Mocks/handlers/authStatusHandlers';
import { server } from '@Mocks/server';
import { screen } from '@testing-library/react';

afterEach(() => {
  server.resetHandlers();
  vi.restoreAllMocks();
});

const renderRouteWithAuthContext = (initialEntry: '/' | '/needAuthentication' | '/needNoAuthentication') =>
  renderRoute({
    initialEntry: initialEntry,
    routes: [
      {
        ...routeObject[0],
        index: false,
        children: [
          {
            path: '/',
            element: <p>home</p>,
          },
          {
            path: '/needAuthentication',
            element: <p>needAuthentication</p>,
            handle: {
              mustBeAuthenticate: true,
            },
          },
          {
            path: '/needNoAuthentication',
            element: <p>needNoAuthentication</p>,
            handle: {
              mustBeAuthenticate: false,
            },
          },

          {
            path: '/sign-up',
            element: <p>signUp</p>,
            handle: {
              mustBeAuthenticate: false,
            },
          },
        ],
      },
    ],
  });

describe('AuthenticationProtection with valid credentials', () => {
  beforeEach(() => {
    server.use(getAuthStatus200);
  });

  it('should render children when route does not require authentication and the user is authenticated', async () => {
    renderRouteWithAuthContext('/');

    expect(await screen.findByText('home')).toBeInTheDocument();
  });

  it('should render children when route requires authentication and the user is authenticated', async () => {
    renderRouteWithAuthContext('/needAuthentication');

    expect(await screen.findByText('needAuthentication')).toBeInTheDocument();
  });

  it('should redirect by default to home when route requires no authentication and the user is authenticated', async () => {
    renderRouteWithAuthContext('/needNoAuthentication');

    expect(await screen.findByText('home')).toBeInTheDocument();
  });

  it('should redirect to postAuthRedirectPath after authentication and reset it', async () => {
    const mockSetPostAuthRedirectPath = vi.fn();
    const mockResetPostAuthRedirectPath = vi.fn();

    vi.spyOn(useAuthenticationContext, 'useAuthenticationContext').mockReturnValue({
      isAuthenticated: true,
      postAuthRedirectPath: '/',
      setPostAuthRedirectPath: mockSetPostAuthRedirectPath,
      resetPostAuthRedirectPath: mockResetPostAuthRedirectPath,
      checkAuthentication: vi.fn(),
    });

    renderRouteWithAuthContext('/needNoAuthentication');

    expect(await screen.findByText('home')).toBeInTheDocument();
    expect(mockResetPostAuthRedirectPath).toHaveBeenCalled();
  });
});

describe('AuthenticationProtection with invalid credentials', () => {
  beforeEach(() => {
    server.use(getAuthStatus400);
  });

  it('should render children when route does not require authentication and the user is not authenticated', async () => {
    renderRouteWithAuthContext('/');

    expect(await screen.findByText('home')).toBeInTheDocument();
  });

  it('should redirect to signUp when route requires authentication and the user is not authenticated', async () => {
    renderRouteWithAuthContext('/needAuthentication');

    expect(await screen.findByText('signUp')).toBeInTheDocument();
  });

  it('should render children when route requires no authentication and the user is not authenticated', async () => {
    renderRouteWithAuthContext('/needNoAuthentication');

    expect(await screen.findByText('needNoAuthentication')).toBeInTheDocument();
  });

  it('should set postAuthRedirectPath when trying to access a protected route while not authenticated', () => {
    const mockSetPostAuthRedirectPath = vi.fn();

    vi.spyOn(useAuthenticationContext, 'useAuthenticationContext').mockReturnValue({
      isAuthenticated: false,
      postAuthRedirectPath: undefined,
      setPostAuthRedirectPath: mockSetPostAuthRedirectPath,
      resetPostAuthRedirectPath: vi.fn(),
      checkAuthentication: vi.fn(),
    });
    renderRouteWithAuthContext('/needAuthentication');

    expect(mockSetPostAuthRedirectPath).toHaveBeenCalledWith('/needAuthentication');
  });
});
