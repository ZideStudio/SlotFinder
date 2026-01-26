import { errorRoutes } from '@Front/pages/Error';
import { homeRoutes } from '@Front/pages/Home';
import { appRoutes } from '@Front/routing/appRoutes';
import { renderRoute } from '@Front/utils/testsUtils/customRender/customRender';
import { screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { oauthCallbackRoutes } from '../routes';

const renderOAuthCallback = (params?: { error?: string; returnUrl?: string }) =>
  renderRoute({
    routes: [oauthCallbackRoutes, homeRoutes, errorRoutes],
    initialEntry: appRoutes.oAuthCallback(params),
  });

describe('OAuthCallback', () => {
  it('should redirect to /error with state when error param is present', async () => {
    renderOAuthCallback({ error: 'TestError' });
    expect(await screen.findByText('TestError')).toBeInTheDocument();
  });

  it('should redirect to / when no error param is present', async () => {
    renderOAuthCallback();
    expect(await screen.findByText('dashboard.title')).toBeInTheDocument();
  });

  it('should redirect to returnUrl when returnUrl param is present', async () => {
    renderOAuthCallback({ returnUrl: appRoutes.home() });
    expect(await screen.findByText('dashboard.title')).toBeInTheDocument();
  });

  it('should redirect to home when returnUrl param is invalid', async () => {
    renderOAuthCallback({ returnUrl: 'invalid-url' });
    expect(await screen.findByText('dashboard.title')).toBeInTheDocument();
  });

  describe('Security: Open Redirect Prevention', () => {
    it('should redirect to home when returnUrl is a protocol-relative URL', async () => {
      renderOAuthCallback({ returnUrl: '//evil.com' });
      expect(await screen.findByText('dashboard.title')).toBeInTheDocument();
    });

    it('should redirect to home when returnUrl is an absolute http URL', async () => {
      renderOAuthCallback({ returnUrl: 'http://evil.com' });
      expect(await screen.findByText('dashboard.title')).toBeInTheDocument();
    });

    it('should redirect to home when returnUrl is an absolute https URL', async () => {
      renderOAuthCallback({ returnUrl: 'https://evil.com' });
      expect(await screen.findByText('dashboard.title')).toBeInTheDocument();
    });

    it('should redirect to home when returnUrl uses javascript protocol', async () => {
      renderOAuthCallback({ returnUrl: 'javascript:alert(1)' });
      expect(await screen.findByText('dashboard.title')).toBeInTheDocument();
    });
  });
});
