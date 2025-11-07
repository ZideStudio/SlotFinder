import { AuthenticationContextProvider } from '@Front/contexts/AuthenticationContext/AuthenticationContextProvider';
import * as useAuthenticationContext from '@Front/hooks/useAuthenticationContext';
import { renderWithQueryClient } from '@Front/utils/testsUtils/customRender/customRender';
import { screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { oauthProvidersData } from '../constants';
import { OAuth } from '../OAuth';

const renderOAuth = () =>
  renderWithQueryClient(
    <AuthenticationContextProvider>
      <OAuth />
    </AuthenticationContextProvider>,
  );

describe('OAuth', () => {
  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('should render heading with correct text and aria-labelledby', () => {
    renderOAuth();
    expect(screen.getByRole('heading', { level: 2, name: 'authentication.signInWithProvider' })).toBeInTheDocument();
  });

  it('should render all OAuth providers as links with correct aria-labels and generated URLs', () => {
    const RETURN_URL = encodeURIComponent('/dashboard');

    renderOAuth();

    for (const provider of oauthProvidersData) {
      const link = screen.getByRole('link', { name: `authentication.${provider.ariaLabel}` });
      expect(link).toBeInTheDocument();
      expect(link).toHaveAttribute(
        'href',
        `${import.meta.env.FRONT_BACKEND_URL}/v1/auth/${provider.id}/url?returnUrl=${RETURN_URL}`,
      );
      expect(link).toHaveAttribute('rel', 'noopener noreferrer');
      expect(screen.getByText(provider.label)).toBeInTheDocument();
    }
  });

  it('should render all OAuth providers as links with correct generated URLs from postAuthRedirectPath', () => {
    const CUSTOM_RETURN_URL = encodeURIComponent('/custom-path');

    vi.spyOn(useAuthenticationContext, 'useAuthenticationContext').mockReturnValue({
      isAuthenticated: false,
      postAuthRedirectPath: '/custom-path',
      setPostAuthRedirectPath: vi.fn(),
      resetPostAuthRedirectPath: vi.fn(),
      checkAuthentication: vi.fn(),
    });

    renderOAuth();

    for (const provider of oauthProvidersData) {
      const link = screen.getByRole('link', { name: `authentication.${provider.ariaLabel}` });
      expect(link).toBeInTheDocument();
      expect(link).toHaveAttribute(
        'href',
        `${import.meta.env.FRONT_BACKEND_URL}/v1/auth/${provider.id}/url?returnUrl=${CUSTOM_RETURN_URL}`,
      );
    }
  });

  it('should render provider icons with aria-hidden', () => {
    renderOAuth();

    for (const provider of oauthProvidersData) {
      const link = screen.getByRole('link', { name: `authentication.${provider.ariaLabel}` });
      const icon = link.querySelector('svg');
      expect(icon).toHaveAttribute('aria-hidden');
    }
  });
});
