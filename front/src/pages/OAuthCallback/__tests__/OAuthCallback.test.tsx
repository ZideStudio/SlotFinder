import { dashboardRoutes } from '@Front/pages/Dashboard';
import { errorRoutes } from '@Front/pages/Error';
import { homeRoutes } from '@Front/pages/Home';
import { appRoutes } from '@Front/routing/appRoutes';
import { renderRoute } from '@Front/utils/testsUtils/customRender/customRender';
import { screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { oauthCallbackRoutes } from '../routes';

const renderOAuthCallback = (params?: { message?: string; returnUrl?: string }) =>
  renderRoute({
    routes: [oauthCallbackRoutes, homeRoutes, dashboardRoutes, errorRoutes],
    initialEntry: appRoutes.oAuthCallback(params),
  });

describe('OAuthCallback', () => {
  it('should redirect to /error with state when message param is present', async () => {
    renderOAuthCallback({ message: 'TestError' });
    expect(await screen.findByText('TestError')).toBeInTheDocument();
  });

  it('should redirect to / when no message param is present', async () => {
    renderOAuthCallback();
    expect(await screen.findByText('dashboard.title')).toBeInTheDocument();
  });

  it('should redirect to returnUrl when returnUrl param is present', async () => {
    renderOAuthCallback({ returnUrl: appRoutes.home() });
    expect(await screen.findByText('home.welcome')).toBeInTheDocument();
  });

  it('should redirect to dashboard when returnUrl param is invalid', async () => {
    renderOAuthCallback({ returnUrl: 'invalid-url' });
    expect(await screen.findByText('dashboard.title')).toBeInTheDocument();
  });
});
