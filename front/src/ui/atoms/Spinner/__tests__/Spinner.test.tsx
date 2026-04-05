import { render, screen } from '@testing-library/react';
import { Spinner } from '../Spinner';

describe('Spinner', () => {
  it('applies the default class name', () => {
    render(<Spinner />);

    const spinner = screen.getByRole('status');
    expect(spinner).toHaveClass('ds-spinner');
    expect(spinner).toHaveAttribute('aria-label', 'Chargement en cours');
  });

  it('applies a custom className when provided', () => {
    render(<Spinner className="my-custom-class" />);

    const spinner = screen.getByRole('status');
    expect(spinner).toHaveClass('ds-spinner', 'my-custom-class');
  });
});
