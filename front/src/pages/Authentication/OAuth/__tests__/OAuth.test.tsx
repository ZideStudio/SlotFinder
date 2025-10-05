import { renderWithQueryClient } from '@Front/utils/testsUtils/customRender';
import { oAuthProvidersErrorFixture, oAuthProvidersFixture } from '@Mocks/fixtures/oAuthProvidersFixtures';
import { getOAuthProviders200, getOAuthProviders400 } from '@Mocks/handlers/oAuthProvidersHandlers';
import { server } from '@Mocks/server';
import { screen, waitFor } from '@testing-library/react';
import { oauthProvidersData } from '../constants';
import { OAuth } from '../index';

beforeAll(() => {
  server.resetHandlers();
});

afterEach(() => {
  server.resetHandlers();
});

const renderOAuth = () => renderWithQueryClient(<OAuth />);

describe('OAuth', () => {
  beforeEach(() => {
    server.use(getOAuthProviders200);
  });

  it('renders heading with correct text and aria-labelledby', () => {
    renderOAuth();
    expect(screen.getByRole('heading', { level: 2, name: 'authentication.signInWithProvider' })).toBeInTheDocument();
  });

  it('renders all OAuth providers as links with correct aria-labels', async () => {
    renderOAuth();

    for (const provider of oauthProvidersData) {
      const link = screen.getByRole('link', { name: `authentication.${provider.ariaLabel}` });
      expect(link).toBeInTheDocument();
      await waitFor(() => expect(link).toHaveAttribute('href', oAuthProvidersFixture[provider.id]));
      expect(link).toHaveAttribute('rel', 'noopener noreferrer');
      expect(screen.getByText(provider.label)).toBeInTheDocument();
    }
  });

  it('renders provider icons with aria-hidden', () => {
    renderOAuth();

    oauthProvidersData.forEach(provider => {
      const link = screen.getByRole('link', { name: `authentication.${provider.ariaLabel}` });
      const icon = link.querySelector('svg');
      expect(icon).toHaveAttribute('aria-hidden');
    });
  });
});

describe('OAuth error handling', () => {
  beforeEach(() => {
    server.use(getOAuthProviders400);
  });

  it('displays an error message when provider link fails', async () => {
    renderOAuth();

    oauthProvidersData.forEach(provider => {
      const link = screen.getByRole('link', { name: `authentication.${provider.ariaLabel}` });
      expect(link).toHaveAttribute('href', '#');
      expect(link).toHaveAttribute('rel', 'noopener noreferrer');
    });

    expect(await screen.findByText(`authentication.error.${oAuthProvidersErrorFixture.code}`)).toBeInTheDocument();
  });
});
