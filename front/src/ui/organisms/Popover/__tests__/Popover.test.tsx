import { render, screen } from '@testing-library/react';

import { Popover } from '../Popover';

describe('Popover', () => {
  describe('rendering', () => {
    it('should render title, content, and footer actions', () => {
      render(
        <Popover
          id="my-popover"
          title="Popover title"
          primaryButtonProps={{ children: 'Confirm' }}
          secondaryButtonProps={{ children: 'Cancel' }}
        >
          Popover content
        </Popover>,
      );

      expect(screen.getByRole('heading', { name: 'Popover title', hidden: true })).toBeInTheDocument();
      expect(screen.getByText('Popover content')).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Confirm', hidden: true })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Cancel', hidden: true })).toBeInTheDocument();
    });

    it('should not render secondary button when secondaryButtonProps is not provided', () => {
      render(
        <Popover id="my-popover" title="Popover title" primaryButtonProps={{ children: 'Confirm' }}>
          Popover content
        </Popover>,
      );

      expect(screen.queryByRole('button', { name: 'Cancel', hidden: true })).not.toBeInTheDocument();
    });

    it('should render a close button with an accessible name', () => {
      render(
        <Popover id="my-popover" title="Popover title" primaryButtonProps={{ children: 'Confirm' }}>
          Popover content
        </Popover>,
      );

      const closeButton = screen.getByRole('button', { name: 'Fermer la fenêtre', hidden: true });
      expect(closeButton).toBeInTheDocument();
      expect(closeButton).toHaveAccessibleName('Fermer la fenêtre');
    });
  });

  describe('popover attribute', () => {
    it('should have popover="auto" attribute by default', () => {
      const { container } = render(
        <Popover id="my-popover" title="Popover title" primaryButtonProps={{ children: 'Confirm' }}>
          Popover content
        </Popover>,
      );

      expect(container.firstChild).toHaveAttribute('popover', 'auto');
    });

    it('should allow overriding the popover attribute', () => {
      const { container } = render(
        <Popover id="my-popover" title="Popover title" primaryButtonProps={{ children: 'Confirm' }} popover="manual">
          Popover content
        </Popover>,
      );

      expect(container.firstChild).toHaveAttribute('popover', 'manual');
    });
  });

  describe('id', () => {
    it('should set the id attribute on the root element', () => {
      const { container } = render(
        <Popover id="my-popover" title="Popover title" primaryButtonProps={{ children: 'Confirm' }}>
          Popover content
        </Popover>,
      );

      expect(container.firstChild).toHaveAttribute('id', 'my-popover');
    });
  });

  describe('close button link', () => {
    it('should link the close button to the popover via popovertarget and popovertargetaction', () => {
      render(
        <Popover id="my-popover" title="Popover title" primaryButtonProps={{ children: 'Confirm' }}>
          Popover content
        </Popover>,
      );

      const closeButton = screen.getByRole('button', { name: 'Fermer la fenêtre', hidden: true });
      expect(closeButton).toHaveAttribute('popovertarget', 'my-popover');
      expect(closeButton).toHaveAttribute('popovertargetaction', 'hide');
    });
  });

  describe('className', () => {
    it('should apply the default ds-popover class', () => {
      const { container } = render(
        <Popover id="my-popover" title="Popover title" primaryButtonProps={{ children: 'Confirm' }}>
          Popover content
        </Popover>,
      );

      expect(container.firstChild).toHaveClass('ds-popover');
      expect(container.firstChild).toHaveClass('ds-card');
    });

    it('should apply a custom className alongside the default class', () => {
      const { container } = render(
        <Popover
          id="my-popover"
          title="Popover title"
          primaryButtonProps={{ children: 'Confirm' }}
          className="custom-popover"
        >
          Popover content
        </Popover>,
      );

      expect(container.firstChild).toHaveClass('ds-popover');
      expect(container.firstChild).toHaveClass('ds-card');
      expect(container.firstChild).toHaveClass('custom-popover');
    });
  });
});
