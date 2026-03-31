import { render, screen } from '@testing-library/react';
import { Spinner } from '../Spinner';

describe('Spinner', () => {
  it('renders the accessible container with correct aria attributes', () => {
    render(<Spinner />);

    const container = screen.getByLabelText('Chargement en cours...');

    expect(container).toHaveAttribute('aria-live', 'polite');
    expect(container).toHaveAttribute('aria-busy', 'true');
  });

  it('renders the spinner div with presentation role and aria-hidden', () => {
    render(<Spinner />);

    const spinner = screen.getByRole('presentation', { hidden: true });

    expect(spinner).toBeInTheDocument();
    expect(spinner).toHaveAttribute('aria-hidden', 'true');
  });

  it('applies the default class name', () => {
    render(<Spinner />);

    const spinner = screen.getByRole('presentation', { hidden: true });
    expect(spinner).toHaveClass('ds-spinner');
  });

  it('applies a custom className when provided', () => {
    render(<Spinner className="my-custom-class" />);

    const spinner = screen.getByRole('presentation', { hidden: true });
    expect(spinner).toHaveClass('ds-spinner', 'my-custom-class');
  });
});
