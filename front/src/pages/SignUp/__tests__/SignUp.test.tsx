import { serverUsePost } from '@Front/utils/testsUtils/msw';
import { server } from '@Mocks/server';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { render, screen } from '@testing-library/react';
import { userEvent } from '@testing-library/user-event';
import { afterAll, afterEach, beforeAll, describe, expect, it } from 'vitest';
import { SignUp } from '../SignUp';

const queryClient = new QueryClient();
const renderWithProvider = (ui: React.ReactElement) =>
  render(<QueryClientProvider client={queryClient}>{ui}</QueryClientProvider>);

beforeAll(() => server.listen());
afterAll(() => server.close());
afterEach(() => server.resetHandlers());

describe('SignUp', () => {
  it('renders all form fields and submit button', () => {
    renderWithProvider(<SignUp />);
    expect(screen.getByLabelText('Username')).toBeInTheDocument();
    expect(screen.getByLabelText('Email')).toBeInTheDocument();
    expect(screen.getByLabelText('Password')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /sign up/i })).toBeInTheDocument();
  });

  it('shows validation errors for empty fields', async () => {
    renderWithProvider(<SignUp />);
    userEvent.click(screen.getByRole('button', { name: /sign up/i }));

    expect(await screen.findByText('Username is required')).toBeInTheDocument();
    expect(screen.getByText('Email is required')).toBeInTheDocument();
    expect(screen.getByText('Password is required')).toBeInTheDocument();
  });

  it.each([
    {
      password: '1234567',
      expectedError: 'Password must be at least 8 characters',
      description: 'minimum length error',
    },
    {
      password: '12345678!',
      expectedError: 'Password must contain letters',
      description: 'must contain letters error',
    },
    {
      password: 'Password!',
      expectedError: 'Password must contain numbers',
      description: 'must contain numbers error',
    },
    {
      password: 'Password1',
      expectedError: 'Password must contain symbols',
      description: 'must contain symbols error',
    },
  ])('shows password $description', async ({ password, expectedError }) => {
    renderWithProvider(<SignUp />);
    await userEvent.type(screen.getByLabelText('Username'), 'testuser');
    await userEvent.type(screen.getByLabelText('Email'), 'test@example.com');
    await userEvent.type(screen.getByLabelText('Password'), password);
    userEvent.click(screen.getByRole('button', { name: 'Sign Up' }));

    expect(await screen.findByText(expectedError)).toBeInTheDocument();
  });

  it('shows error message on failed submission', async () => {
    serverUsePost({
      route: '/v1/account',
      code: 400,
      responseBody: { error: true, code: 'UNKNOWN_ERROR', message: 'Username already exists' },
    });

    renderWithProvider(<SignUp />);
    await userEvent.type(screen.getByLabelText('Username'), 'failuser');
    await userEvent.type(screen.getByLabelText('Email'), 'fail@example.com');
    await userEvent.type(screen.getByLabelText('Password'), 'Password1!');
    userEvent.click(screen.getByRole('button', { name: /sign up/i }));

    expect(await screen.findByText('Username already exists')).toBeInTheDocument();
  });
});
