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
    await expect(screen.findByText('TestError')).resolves.toBeInTheDocument();
  });

  it('should redirect to / when no error param is present', async () => {
    renderOAuthCallback();
    await expect(screen.findByText('dashboard.title')).resolves.toBeInTheDocument();
  });

  it('should redirect to returnUrl when returnUrl param is present', async () => {
    renderOAuthCallback({ returnUrl: appRoutes.home() });
    await expect(screen.findByText('dashboard.title')).resolves.toBeInTheDocument();
  });

  it('should redirect to home when returnUrl param is invalid', async () => {
    renderOAuthCallback({ returnUrl: 'invalid-url' });
    await expect(screen.findByText('dashboard.title')).resolves.toBeInTheDocument();
  });

  describe('Security: Open Redirect Prevention', () => {
    it('should redirect to home when returnUrl is a protocol-relative URL', async () => {
      renderOAuthCallback({ returnUrl: '//evil.com' });
      await expect(screen.findByText('dashboard.title')).resolves.toBeInTheDocument();
    });

    it('should redirect to home when returnUrl is an absolute http URL', async () => {
      renderOAuthCallback({ returnUrl: 'http://evil.com' });
      await expect(screen.findByText('dashboard.title')).resolves.toBeInTheDocument();
    });

    it('should redirect to home when returnUrl is an absolute https URL', async () => {
      renderOAuthCallback({ returnUrl: 'https://evil.com' });
      await expect(screen.findByText('dashboard.title')).resolves.toBeInTheDocument();
    });

    it('should redirect to home when returnUrl uses javascript protocol', async () => {
      // oxlint-disable-next-line no-script-url
      renderOAuthCallback({ returnUrl: 'javascript:alert(1)' });
      await expect(screen.findByText('dashboard.title')).resolves.toBeInTheDocument();
    });
  });
});
