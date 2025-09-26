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
  vi.mock('react-i18next', () => ({
    useTranslation: vi.fn().mockReturnValue({
      t: (messageId: string, args: Record<string, unknown>) => messageId + (args ? `::${JSON.stringify(args)}` : ''),
    }),
  }));

  it('renders all form fields and submit button', () => {
    renderWithProvider(<SignUp />);
    expect(screen.getByLabelText('username')).toBeInTheDocument();
    expect(screen.getByLabelText('email')).toBeInTheDocument();
    expect(screen.getByLabelText('password')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'submit' })).toBeInTheDocument();
  });

  it('shows validation errors for empty fields', async () => {
    renderWithProvider(<SignUp />);
    await userEvent.click(screen.getByRole('button', { name: 'submit' }));

    expect(await screen.findByText('requiredUsername')).toBeInTheDocument();
    expect(screen.getByText('requiredEmail')).toBeInTheDocument();
    expect(screen.getByText('requiredPassword')).toBeInTheDocument();
  });

  it.each([
    {
      password: '1234567',
      expectedError: 'minLengthPassword::{"min":8}',
      description: 'minimum length error',
    },
    {
      password: '12345678!',
      expectedError: 'passwordComplexity',
      description: 'must contain letters error',
    },
    {
      password: 'Password!',
      expectedError: 'passwordComplexity',
      description: 'must contain numbers error',
    },
    {
      password: 'Password1',
      expectedError: 'passwordComplexity',
      description: 'must contain symbols error',
    },
  ])('shows password $description', async ({ password, expectedError }) => {
    renderWithProvider(<SignUp />);
    await userEvent.type(screen.getByLabelText('username'), 'testuser');
    await userEvent.type(screen.getByLabelText('email'), 'test@example.com');
    await userEvent.type(screen.getByLabelText('password'), password);
    await userEvent.click(screen.getByRole('button', { name: 'submit' }));

    expect(await screen.findByText(expectedError)).toBeInTheDocument();
  });

  it('shows error message on failed submission', async () => {
    serverUsePost({
      route: '/v1/account',
      code: 400,
      responseBody: { error: true, code: 'UNKNOWN_ERROR', message: 'Username already exists' },
    });

    renderWithProvider(<SignUp />);
    await userEvent.type(screen.getByLabelText('username'), 'failuser');
    await userEvent.type(screen.getByLabelText('email'), 'fail@example.com');
    await userEvent.type(screen.getByLabelText('password'), 'Password1!');
    await userEvent.click(screen.getByRole('button', { name: 'submit' }));

    expect(await screen.findByText('Username already exists')).toBeInTheDocument();
  });
});
