import { render, screen } from '@testing-library/react';
import { oauthProviders } from '../constants';
import { OAuth } from '../index';

describe('OAuth', () => {
  it('renders heading with correct text and aria-labelledby', () => {
    render(<OAuth />);
    expect(screen.getByRole('heading', { level: 2, name: 'signInWithProvider' })).toBeInTheDocument();
  });

  it('renders all OAuth providers as links with correct aria-labels', () => {
    render(<OAuth />);

    oauthProviders.forEach(provider => {
      const link = screen.getByRole('link', { name: provider.ariaLabel });
      expect(link).toBeInTheDocument();
      expect(link).toHaveAttribute('href', provider.href);
      expect(link).toHaveAttribute('rel', 'noopener noreferrer');
      expect(screen.getByText(provider.label)).toBeInTheDocument();
    });
  });

  it('renders provider icons with aria-hidden', () => {
    render(<OAuth />);

    oauthProviders.forEach(provider => {
      const link = screen.getByRole('link', { name: provider.ariaLabel });
      const icon = link.querySelector('svg');
      expect(icon).toHaveAttribute('aria-hidden');
    });
  });
});
