import { render, fireEvent, screen } from '@testing-library/react';
import { Toast } from '../Toast';

describe('Toast', () => {
  it('renders the Toast component with the provided message', () => {
    const message = 'This is a toast message';
    render(<Toast onClose={() => {}}>{message}</Toast>);
    const toastElement = screen.getByText(message);
    expect(toastElement).toBeInTheDocument();
  });

  it('calls the onClose function when the close button is clicked', () => {
    const onCloseMock = vitest.fn();
    render(<Toast onClose={onCloseMock}>Test Toast</Toast>);
    const closeButton = screen.getByRole('button');
    fireEvent.click(closeButton);
    expect(onCloseMock).toHaveBeenCalledTimes(1);
  });

  it('applies the correct class names', () => {
    render(
      <Toast onClose={() => {}} className="custom-toast">
        Test Toast
      </Toast>,
    );
    const toastElement = screen.getByRole('status');
    expect(toastElement).toHaveClass('ds-toast ds-toast--custom-toast');
  });

  it('is visible after mounting', () => {
    render(<Toast onClose={() => {}}>Test Toast</Toast>);
    const toastElement = screen.getByRole('status');
    setTimeout(() => {
      expect(toastElement).toHaveClass('ds-toast--visible');
    }, 100);
  });
});
