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

  it('applies a custom aria-label when provided', () => {
    render(<Spinner aria-label="Chargement des données" />);
    const spinner = screen.getByRole('status', {
      name: 'Chargement des données',
    });
    expect(spinner).toHaveAttribute('aria-label', 'Chargement des données');
  });

  it('applies the size style when size prop is provided', () => {
    render(<Spinner size="32px" />);
    const spinner = screen.getByRole('status');
    expect(spinner).toHaveStyle('--ds-spinner-size: 32px');
  });
});
