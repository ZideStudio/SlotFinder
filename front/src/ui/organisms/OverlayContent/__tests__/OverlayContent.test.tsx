import { render, screen, within } from '@testing-library/react';
import { userEvent } from '@testing-library/user-event';

import { OverlayContent } from '../OverlayContent';

describe('OverlayContent', () => {
  describe('rendering', () => {
    it('should render title as a level 1 heading, content, and footer actions', () => {
      render(
        <OverlayContent
          title="Content title"
          primaryButtonProps={{ children: 'Confirm' }}
          secondaryButtonProps={{ children: 'Cancel' }}
          closeOverlay={vi.fn()}
        >
          Body content
        </OverlayContent>,
      );

      expect(screen.getByRole('heading', { name: 'Content title', level: 1 })).toBeInTheDocument();
      expect(screen.getByText('Body content')).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Confirm' })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Cancel' })).toBeInTheDocument();
    });

    it('should not render secondary button when secondaryButtonProps is not provided', () => {
      render(
        <OverlayContent title="Content title" primaryButtonProps={{ children: 'Confirm' }} closeOverlay={vi.fn()}>
          Body content
        </OverlayContent>,
      );

      expect(screen.getByRole('button', { name: 'Confirm' })).toBeInTheDocument();
      expect(screen.queryByRole('button', { name: 'Cancel' })).not.toBeInTheDocument();
    });

    it('should render a close button with accessible name and structural sections', () => {
      render(
        <OverlayContent title="Content title" primaryButtonProps={{ children: 'Confirm' }} closeOverlay={vi.fn()}>
          Body content
        </OverlayContent>,
      );

      const closeButton = screen.getByRole('button', { name: 'Fermer la fenêtre' });
      expect(closeButton).toBeInTheDocument();
      expect(closeButton).toHaveAccessibleName('Fermer la fenêtre');

      expect(screen.getByRole('banner')).toBeInTheDocument();
      expect(screen.getByText('Body content')).toBeInTheDocument();

      const footer = screen.getByRole('contentinfo');
      expect(footer).toBeInTheDocument();
      expect(within(footer).getByRole('button', { name: 'Confirm' })).toBeInTheDocument();
    });
  });

  describe('titleId', () => {
    it('should set the id attribute on the heading when titleId is provided', () => {
      render(
        <OverlayContent
          title="Content title"
          titleId="custom-title-id"
          primaryButtonProps={{ children: 'Confirm' }}
          closeOverlay={vi.fn()}
        >
          Body content
        </OverlayContent>,
      );

      expect(screen.getByRole('heading', { name: 'Content title' })).toHaveAttribute('id', 'custom-title-id');
    });

    it('should not set the id attribute on the heading when titleId is not provided', () => {
      render(
        <OverlayContent title="Content title" primaryButtonProps={{ children: 'Confirm' }} closeOverlay={vi.fn()}>
          Body content
        </OverlayContent>,
      );

      expect(screen.getByRole('heading', { name: 'Content title' })).not.toHaveAttribute('id');
    });
  });

  describe('interactions', () => {
    it('should call closeOverlay when the close button is clicked', async () => {
      const closeOverlay = vi.fn();

      render(
        <OverlayContent title="Content title" primaryButtonProps={{ children: 'Confirm' }} closeOverlay={closeOverlay}>
          Body content
        </OverlayContent>,
      );

      await userEvent.click(screen.getByRole('button', { name: 'Fermer la fenêtre' }));
      expect(closeOverlay).toHaveBeenCalledTimes(1);
    });

    it('should call primaryButtonProps onClick when primary button is clicked', async () => {
      const onClick = vi.fn();

      render(
        <OverlayContent
          title="Content title"
          primaryButtonProps={{ children: 'Confirm', onClick }}
          closeOverlay={vi.fn()}
        >
          Body content
        </OverlayContent>,
      );

      await userEvent.click(screen.getByRole('button', { name: 'Confirm' }));
      expect(onClick).toHaveBeenCalledTimes(1);
    });

    it('should call secondaryButtonProps onClick when secondary button is clicked', async () => {
      const onClick = vi.fn();

      render(
        <OverlayContent
          title="Content title"
          primaryButtonProps={{ children: 'Confirm' }}
          secondaryButtonProps={{ children: 'Cancel', onClick }}
          closeOverlay={vi.fn()}
        >
          Body content
        </OverlayContent>,
      );

      await userEvent.click(screen.getByRole('button', { name: 'Cancel' }));
      expect(onClick).toHaveBeenCalledTimes(1);
    });
  });
});
