import { appRoutes } from '@Front/routing/appRoutes';
import { renderRoute, type RenderRouteOptions } from '@Front/utils/testsUtils/customRender';
import { accountErrorFixture } from '@Mocks/fixtures/accountFixtures';
import { postAccount201, postAccount400 } from '@Mocks/handlers/accountHandlers';
import { getOAuthProviders200 } from '@Mocks/handlers/oAuthProvidersHandlers';
import { server } from '@Mocks/server';
import { screen } from '@testing-library/react';
import { userEvent } from '@testing-library/user-event';
import { describe, expect, it } from 'vitest';
import { authenticationRoutes } from '../../routes';

const renderRouteOptions: RenderRouteOptions = {
  routes: [authenticationRoutes],
  routesOptions: { initialEntries: [appRoutes.signUp()] },
};

beforeEach(() => {
  server.use(getOAuthProviders200);
});

afterEach(() => {
  server.resetHandlers();
});

describe('SignUp', () => {
  beforeEach(() => {
    server.use(postAccount201);
  });

  it('renders all form fields, submit button and oauth', () => {
    renderRoute(renderRouteOptions);

    expect(screen.getByLabelText('signUp.username')).toBeInTheDocument();
    expect(screen.getByLabelText('signUp.email')).toBeInTheDocument();
    expect(screen.getByLabelText('signUp.password')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'signUp.submit' })).toBeInTheDocument();
    expect(screen.getByRole('heading', { level: 2, name: 'authentication.signInWithProvider' })).toBeInTheDocument();
  });

  it('shows validation errors for empty fields', async () => {
    renderRoute(renderRouteOptions);

    await userEvent.click(screen.getByRole('button', { name: 'signUp.submit' }));

    expect(await screen.findByText('signUp.requiredUsername')).toBeInTheDocument();
    expect(screen.getByText('signUp.requiredEmail')).toBeInTheDocument();
    expect(screen.getByText('signUp.requiredPassword')).toBeInTheDocument();
  });

  it.each([
    {
      password: '1234567',
      expectedError: 'signUp.minLengthPassword::{"min":8}',
      description: 'minimum length error',
    },
    {
      password: '12345678!',
      expectedError: 'signUp.passwordComplexity',
      description: 'must contain letters error',
    },
    {
      password: 'password1!',
      expectedError: 'signUp.passwordComplexity',
      description: 'must contain numbers error',
    },
    {
      password: 'Password!',
      expectedError: 'signUp.passwordComplexity',
      description: 'must contain numbers error',
    },
    {
      password: 'Password1',
      expectedError: 'signUp.passwordComplexity',
      description: 'must contain symbols error',
    },
  ])('shows password $description', async ({ password, expectedError }) => {
    renderRoute(renderRouteOptions);

    await userEvent.type(screen.getByLabelText('signUp.username'), 'testuser');
    await userEvent.type(screen.getByLabelText('signUp.email'), 'test@example.com');
    await userEvent.type(screen.getByLabelText('signUp.password'), password);
    await userEvent.click(screen.getByRole('button', { name: 'signUp.submit' }));

    expect(await screen.findByText(expectedError)).toBeInTheDocument();
  });
});

describe('SignUp error handling', () => {
  beforeEach(() => {
    server.use(postAccount400);
  });

  it('shows error message on failed submission', async () => {
    renderRoute(renderRouteOptions);

    await userEvent.type(screen.getByLabelText('signUp.username'), 'failuser');
    await userEvent.type(screen.getByLabelText('signUp.email'), 'fail@example.com');
    await userEvent.type(screen.getByLabelText('signUp.password'), 'Password1!');
    await userEvent.click(screen.getByRole('button', { name: 'signUp.submit' }));

    expect(await screen.findByText(`signUp.error.${accountErrorFixture.code}`)).toBeInTheDocument();
  });
});
