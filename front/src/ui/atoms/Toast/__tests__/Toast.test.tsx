import { render, fireEvent, screen } from '@testing-library/react';
import { useToastService } from '@Front/ui/utils/toast/hooks/useToastService';
import { ToastProvider } from '@Front/ui/utils/toast/toastProvider/ToastProvider';

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
