import { render, fireEvent, act } from '@testing-library/react';
import { ToastProvider } from '../ToastProvider';
import { useToastService } from '@Front/hooks/useToastService';
import { screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

describe('ToastProvider', () => {
  afterEach(() => {
    vi.useRealTimers();
  });

  const TestComponent = () => {
    const toast = useToastService();
    return <button onClick={() => toast.addToast('Test Toast')}>Show Toast</button>;
  };

  it('should render children and provide toast context', async () => {
    render(
      <ToastProvider>
        <TestComponent />
      </ToastProvider>,
    );

    const button = screen.getByRole('button', { name: 'Show Toast' });
    expect(button).toBeInTheDocument();

    await userEvent.click(button);
    const toastMessage = await screen.findByText('Test Toast');

    expect(toastMessage).toBeInTheDocument();
  });

  it('should remove toast after duration', async () => {
    vi.useFakeTimers();

    render(
      <ToastProvider>
        <TestComponent />
      </ToastProvider>,
    );

    fireEvent.click(screen.getByRole('button', { name: 'Show Toast' }));

    expect(screen.getByText('Test Toast')).toBeInTheDocument();

    await act(async () => {
      await vi.runAllTimersAsync();
    });

    expect(screen.queryByText('Test Toast')).not.toBeInTheDocument();
  });

  it('should remove toast when close button is clicked', async () => {
    const { getByText, queryByText } = render(
      <ToastProvider>
        <TestComponent />
      </ToastProvider>,
    );

    const button = getByText('Show Toast');
    await userEvent.click(button);
    const toastMessage = await screen.findByText('Test Toast');
    expect(toastMessage).toBeInTheDocument();

    const closeButton = screen.getByRole('button', { name: 'Fermer la notification' });
    await userEvent.click(closeButton);
    expect(queryByText('Test Toast')).not.toBeInTheDocument();
  });
});
