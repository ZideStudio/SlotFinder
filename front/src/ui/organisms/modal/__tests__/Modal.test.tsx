import { render, screen, within } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

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

  it('should not render secondary action when secondaryButtonProps is not provided', () => {
    render(
      <Modal open title="Modal title" primaryButtonProps={{ children: 'Action' }}>
        Modal content
      </Modal>,
    );

    const dialog = screen.getByRole('dialog', { hidden: true });
    expect(within(dialog).getByRole('button', { name: 'Action' })).toBeInTheDocument();
    expect(within(dialog).queryByRole('button', { name: 'Close' })).not.toBeInTheDocument();
  });

  it('should link dialog label to the title heading and set closedby attribute', () => {
    const { container } = render(
      <Modal open title="Modal title" primaryButtonProps={{ children: 'Action' }}>
        Modal content
      </Modal>,
    );

    const dialog = container.querySelector('dialog');
    const heading = screen.getByRole('heading', { name: 'Modal title' });

    expect(dialog).toBeInTheDocument();
    expect(dialog).toHaveAttribute('aria-labelledby', heading.id);
    expect(dialog).toHaveAttribute('closedby', 'any');
  });

  it('should force secondary button variant to secondary', () => {
    render(
      <Modal
        open
        title="Modal title"
        primaryButtonProps={{ children: 'Action' }}
        secondaryButtonProps={{ children: 'Close', variant: 'primary' }}
      >
        Modal content
      </Modal>,
    );

    expect(screen.getByRole('button', { name: 'Close' })).toHaveClass('ds-button--secondary');
  });

  it('should forward native dialog attributes', () => {
    const onClose = vi.fn();
    const { container } = render(
      <Modal title="Modal title" primaryButtonProps={{ children: 'Action' }} open onClose={onClose}>
        Modal content
      </Modal>,
    );

    const dialog = container.querySelector('dialog');

    expect(dialog).toHaveAttribute('open');
  });
});
