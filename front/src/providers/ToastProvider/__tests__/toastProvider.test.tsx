import { render, fireEvent } from '@testing-library/react';
import { ToastProvider } from '../ToastProvider';
import { useToastContext } from '@Front/hooks/useToastService';

describe('ToastProvider', () => {
  it('should render children and provide toast context', () => {
    const TestComponent = () => {
      const { show } = useToastContext();
      return <button onClick={() => show('Test Toast')}>Show Toast</button>;
    };

    const { getByText } = render(
      <ToastProvider>
        <TestComponent />
      </ToastProvider>,
    );

    const button = getByText('Show Toast');
    expect(button).toBeInTheDocument();

    fireEvent.click(button);
    const toastMessage = getByText('Test Toast');
    expect(toastMessage).toBeInTheDocument();
  });

  it('should remove toast after duration', () => {
    const TestComponent = () => {
      const { show } = useToastContext();
      return <button onClick={() => show('Test Toast')}>Show Toast</button>;
    };

    const { getByText, queryByText } = render(
      <ToastProvider>
        <TestComponent />
      </ToastProvider>,
    );

    const button = getByText('Show Toast');
    fireEvent.click(button);
    const toastMessage = getByText('Test Toast');
    expect(toastMessage).toBeInTheDocument();

    setTimeout(() => {
      expect(queryByText('Test Toast')).not.toBeInTheDocument();
    }, 3500);
  });

  it('should remove toast when close button is clicked', () => {
    const TestComponent = () => {
      const { show } = useToastContext();
      return <button onClick={() => show('Test Toast')}>Show Toast</button>;
    };

    const { getByText, queryByText } = render(
      <ToastProvider>
        <TestComponent />
      </ToastProvider>,
    );

    const button = getByText('Show Toast');
    fireEvent.click(button);
    const toastMessage = getByText('Test Toast');
    expect(toastMessage).toBeInTheDocument();

    const closeButton = getByText('✕');
    fireEvent.click(closeButton);
    expect(queryByText('Test Toast')).not.toBeInTheDocument();
  });
});
