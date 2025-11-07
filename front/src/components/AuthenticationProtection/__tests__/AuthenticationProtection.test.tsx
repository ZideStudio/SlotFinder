import * as useAuthenticationContext from '@Front/hooks/useAuthenticationContext';
import { appRoutes } from '@Front/routing/appRoutes';
import { renderRoute } from '@Front/utils/testsUtils/customRender/customRender';
import { getAuthStatus200, getAuthStatus400 } from '@Mocks/handlers/authStatusHandlers';
import { server } from '@Mocks/server';
import { screen } from '@testing-library/react';

afterEach(() => {
  server.resetHandlers();
  vi.restoreAllMocks();
});

describe('AuthenticationProtection with valid credentials', () => {
  beforeEach(() => {
    server.use(getAuthStatus200);
  });

  it('should render children when route does not require authentication and the user is authenticated', async () => {
    renderRoute({ initialEntry: appRoutes.home() });

    expect(await screen.findByText('home.welcome')).toBeInTheDocument();
  });

  it('should render children when route requires authentication and the user is authenticated', async () => {
    renderRoute({ initialEntry: appRoutes.dashboard() });

    expect(await screen.findByText('dashboard.title')).toBeInTheDocument();
  });

  it('should redirect by default to dashboard when route requires no authentication and the user is authenticated', async () => {
    renderRoute({ initialEntry: appRoutes.signUp() });

    expect(await screen.findByText('dashboard.title')).toBeInTheDocument();
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

    renderRoute({ initialEntry: appRoutes.signUp() });

    expect(await screen.findByText('home.welcome')).toBeInTheDocument();
    expect(mockResetPostAuthRedirectPath).toHaveBeenCalled();
  });
});

describe('AuthenticationProtection with invalid credentials', () => {
  beforeEach(() => {
    server.use(getAuthStatus400);
  });

  it('should render children when route does not require authentication and the user is not authenticated', async () => {
    renderRoute({ initialEntry: appRoutes.home() });

    expect(await screen.findByText('home.welcome')).toBeInTheDocument();
  });

  it('should redirect to signUp when route requires authentication and the user is not authenticated', async () => {
    renderRoute({ initialEntry: appRoutes.dashboard() });

    expect(await screen.findByText('signUp.title')).toBeInTheDocument();
  });

  it('should render children when route requires no authentication and the user is not authenticated', async () => {
    renderRoute({ initialEntry: appRoutes.signUp() });

    expect(await screen.findByText('signUp.title')).toBeInTheDocument();
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
    renderRoute({ initialEntry: appRoutes.dashboard() });

    expect(mockSetPostAuthRedirectPath).toHaveBeenCalledWith('/dashboard');
  });
});
