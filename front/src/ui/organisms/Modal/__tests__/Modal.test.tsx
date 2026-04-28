import { render, screen, within } from '@testing-library/react';
import { userEvent } from '@testing-library/user-event';

import { Modal } from '../Modal';

describe('Modal', () => {
  it('should render title, content, and footer actions', () => {
    render(
      <Modal
        open
        title="Modal title"
        primaryButtonProps={{ children: 'Action' }}
        secondaryButtonProps={{ children: 'Close' }}
      >
        Modal content
      </Modal>,
    );

    const dialog = screen.getByRole('dialog', { hidden: true });
    expect(dialog).toBeInTheDocument();
    within(dialog).getByRole('heading', { name: 'Modal title' });
    expect(within(dialog).getByText('Modal content')).toBeInTheDocument();
    expect(within(dialog).getByRole('button', { name: 'Action' })).toBeInTheDocument();
    expect(within(dialog).getByRole('button', { name: 'Close' })).toBeInTheDocument();
  });

  it('should link dialog label to the title heading and set closedby attribute', () => {
    render(
      <Modal open title="Modal title" primaryButtonProps={{ children: 'Action' }}>
        Modal content
      </Modal>,
    );

    const dialog = screen.getByRole('dialog', { hidden: true });
    const heading = screen.getByRole('heading', { name: 'Modal title' });

    expect(dialog).toBeInTheDocument();
    expect(dialog).toHaveAttribute('aria-labelledby', heading.id);
    expect(dialog).toHaveAttribute('closedby', 'any');
  });

  it('should apply the variant passed via secondaryButtonProps', () => {
    render(
      <Modal
        open
        title="Modal title"
        primaryButtonProps={{ children: 'Action' }}
        secondaryButtonProps={{ children: 'Close', variant: 'secondary' }}
      >
        Modal content
      </Modal>,
    );

    expect(screen.getByRole('button', { name: 'Close' })).toHaveClass('ds-button--secondary');
  });

  it('should forward native dialog attributes', () => {
    const onClose = vi.fn();
    render(
      <Modal title="Modal title" primaryButtonProps={{ children: 'Action' }} open onClose={onClose}>
        Modal content
      </Modal>,
    );

    expect(screen.getByRole('dialog', { hidden: true })).toHaveAttribute('open');
  });

  it('should apply a custom className alongside the default class', () => {
    render(
      <Modal open title="Modal title" primaryButtonProps={{ children: 'Action' }} className="custom-modal">
        Modal content
      </Modal>,
    );

    const dialog = screen.getByRole('dialog', { hidden: true });
    expect(dialog).toHaveClass('ds-modal');
    expect(dialog).toHaveClass('custom-modal');
  });

  it('should call closeModal when the close button is clicked', async () => {
    const close = vi.fn();
    Object.defineProperty(HTMLDialogElement.prototype, 'close', {
      value: close,
      configurable: true,
      writable: true,
    });

    render(
      <Modal open title="Modal title" primaryButtonProps={{ children: 'Action' }}>
        Modal content
      </Modal>,
    );

    await userEvent.click(screen.getByRole('button', { name: 'Fermer la fenêtre' }));
    expect(close).toHaveBeenCalledTimes(1);
  });
});
