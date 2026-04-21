// oxlint-disable-next-line import/no-namespace
import * as useAuthenticationContext from '@Front/hooks/useAuthenticationContext';
import { routeObject } from '@Front/routing/routes';
import { renderRoute } from '@Front/utils/testsUtils/customRender/customRender';
import { getAuthStatus200, getAuthStatus400 } from '@Mocks/handlers/authStatusHandlers';
import { server } from '@Mocks/server';
import { screen } from '@testing-library/react';

describe('AuthenticationProtection', () => {
  type TestRoute = '/' | '/needAuthentication' | '/needNoAuthentication';

  const renderRouteWithAuthContext = (initialEntry: TestRoute) =>
    renderRoute({
      initialEntry,
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

  const mockAuthenticationContext = (
    overrides: Partial<ReturnType<typeof useAuthenticationContext.useAuthenticationContext>> = {},
  ) => {
    const contextValue = {
      isAuthenticated: true,
      postAuthRedirectPath: undefined,
      setPostAuthRedirectPath: vi.fn(),
      resetPostAuthRedirectPath: vi.fn(),
      checkAuthentication: vi.fn(),
      ...overrides,
    } as ReturnType<typeof useAuthenticationContext.useAuthenticationContext>;

    vi.spyOn(useAuthenticationContext, 'useAuthenticationContext').mockReturnValue(contextValue);

    return contextValue;
  };

  const expectDisplayedText = async (initialEntry: TestRoute, text: string) => {
    renderRouteWithAuthContext(initialEntry);

    await expect(screen.findByText(text)).resolves.toBeInTheDocument();
  };

  afterEach(() => {
    server.resetHandlers();
    vi.restoreAllMocks();
  });

  describe('AuthenticationProtection with valid credentials', () => {
    beforeEach(() => {
      server.use(getAuthStatus200);
    });

    it('should render children when route does not require authentication and the user is authenticated', async () => {
      await expectDisplayedText('/', 'home');
    });

    it('should render children when route requires authentication and the user is authenticated', async () => {
      await expectDisplayedText('/needAuthentication', 'needAuthentication');
    });

    it('should redirect by default to home when route requires no authentication and the user is authenticated', async () => {
      await expectDisplayedText('/needNoAuthentication', 'home');
    });

    it('should redirect to postAuthRedirectPath after authentication and reset it', async () => {
      const mockResetPostAuthRedirectPath = vi.fn();

      mockAuthenticationContext({
        isAuthenticated: true,
        postAuthRedirectPath: '/',
        resetPostAuthRedirectPath: mockResetPostAuthRedirectPath,
      });

      await expectDisplayedText('/needNoAuthentication', 'home');
      expect(mockResetPostAuthRedirectPath).toHaveBeenCalledTimes(1);
    });

    it('should not reset postAuthRedirectPath on a neutral route without authentication requirement metadata', async () => {
      const mockResetPostAuthRedirectPath = vi.fn();

      mockAuthenticationContext({
        isAuthenticated: true,
        postAuthRedirectPath: '/',
        resetPostAuthRedirectPath: mockResetPostAuthRedirectPath,
      });

      await expectDisplayedText('/', 'home');
      expect(mockResetPostAuthRedirectPath).not.toHaveBeenCalled();
    });

    it('should not reset postAuthRedirectPath when no redirect path is stored', async () => {
      const mockResetPostAuthRedirectPath = vi.fn();

      mockAuthenticationContext({
        isAuthenticated: true,
        postAuthRedirectPath: undefined,
        resetPostAuthRedirectPath: mockResetPostAuthRedirectPath,
      });

      await expectDisplayedText('/needNoAuthentication', 'home');
      expect(mockResetPostAuthRedirectPath).not.toHaveBeenCalled();
    });
  });

  describe('AuthenticationProtection with pending authentication state', () => {
    it('should render nothing while authentication status is unresolved', () => {
      mockAuthenticationContext({
        isAuthenticated: undefined,
      });

      renderRouteWithAuthContext('/needAuthentication');

      expect(screen.queryByText('needAuthentication')).not.toBeInTheDocument();
      expect(screen.queryByText('signUp')).not.toBeInTheDocument();
      expect(screen.queryByText('home')).not.toBeInTheDocument();
    });
  });

  describe('AuthenticationProtection with invalid credentials', () => {
    beforeEach(() => {
      server.use(getAuthStatus400);
    });

    it('should render children when route does not require authentication and the user is not authenticated', async () => {
      await expectDisplayedText('/', 'home');
    });

    it('should redirect to signUp when route requires authentication and the user is not authenticated', async () => {
      await expectDisplayedText('/needAuthentication', 'signUp');
    });

    it('should render children when route requires no authentication and the user is not authenticated', async () => {
      await expectDisplayedText('/needNoAuthentication', 'needNoAuthentication');
    });

    it('should set postAuthRedirectPath when trying to access a protected route while not authenticated', () => {
      const mockSetPostAuthRedirectPath = vi.fn();

      mockAuthenticationContext({
        isAuthenticated: false,
        postAuthRedirectPath: undefined,
        setPostAuthRedirectPath: mockSetPostAuthRedirectPath,
      });

      renderRouteWithAuthContext('/needAuthentication');

      expect(mockSetPostAuthRedirectPath).toHaveBeenCalledWith('/needAuthentication');
    });
  });
});
