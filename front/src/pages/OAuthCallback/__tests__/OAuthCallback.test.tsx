import { errorRoutes } from '@Front/pages/Error';
import { homeRoutes } from '@Front/pages/Home';
import { appRoutes } from '@Front/routing/appRoutes';
import { renderRoute } from '@Front/utils/testsUtils/customRender';
import { screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { oauthCallbackRoutes } from '../routes';

const renderOAuthCallback = (message?: string) =>
  renderRoute({
    routes: [oauthCallbackRoutes, homeRoutes, errorRoutes],
    routesOptions: { initialEntries: [appRoutes.oAuthCallback(message)] },
  });

describe('OAuthCallback', () => {
  it('should redirect to /error with state when message param is present', async () => {
    renderOAuthCallback('TestError');
    expect(await screen.findByText('TestError')).toBeInTheDocument();
  });

  it('should redirect to / when no message param is present', async () => {
    renderOAuthCallback();
    expect(await screen.findByText('welcome')).toBeInTheDocument();
  });
});
