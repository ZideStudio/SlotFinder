import { renderWithQueryClient } from '@Front/utils/testsUtils/customRender';
import { screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { oauthProvidersData } from '../constants';
import { OAuth } from '../index';

const renderOAuth = () => renderWithQueryClient(<OAuth />);

describe('OAuth', () => {
  it('should render heading with correct text and aria-labelledby', () => {
    renderOAuth();
    expect(screen.getByRole('heading', { level: 2, name: 'authentication.signInWithProvider' })).toBeInTheDocument();
  });

  it('should render all OAuth providers as links with correct aria-labels and generated URLs', () => {
    const REDIRECT_URL = encodeURIComponent('/dashboard');

    renderOAuth();

    for (const provider of oauthProvidersData) {
      const link = screen.getByRole('link', { name: `authentication.${provider.ariaLabel}` });
      expect(link).toBeInTheDocument();
      expect(link).toHaveAttribute(
        'href',
        `${import.meta.env.FRONT_BACKEND_URL}/v1/auth/${provider.id}/url?redirectUrl=${REDIRECT_URL}`,
      );
      expect(link).toHaveAttribute('rel', 'noopener noreferrer');
      expect(screen.getByText(provider.label)).toBeInTheDocument();
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
