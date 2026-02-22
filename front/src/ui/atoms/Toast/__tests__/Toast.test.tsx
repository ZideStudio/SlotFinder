import { render, fireEvent, screen } from '@testing-library/react';
import { useToastService } from '@Front/hooks/useToastService';
import { ToastProvider } from '@Front/providers/ToastProvider/ToastProvider';

describe('Toast', () => {
  it('renders the Toast component with the provided message', () => {
    const message = 'This is a toast message';

    const TestComponent = () => {
      const toast = useToastService();
      return <button onClick={() => toast.addToast(message)}>Show Toast</button>;
    };

    render(
      <ToastProvider>
        <TestComponent />
      </ToastProvider>,
    );

    fireEvent.click(screen.getByText('Show Toast'));
    const toastElement = screen.getByText(message);
    expect(toastElement).toBeInTheDocument();
  });
});
