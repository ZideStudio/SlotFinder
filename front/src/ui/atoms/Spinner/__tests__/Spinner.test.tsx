import { render, screen } from '@testing-library/react';
import { Spinner } from '../Spinner';

describe('Spinner', () => {
  it('applies the default class name', () => {
    render(<Spinner />);

    const spinner = screen.getByRole('status');
    expect(spinner).toHaveClass('ds-spinner');
  });

  it('applies a custom className when provided', () => {
    render(<Spinner className="my-custom-class" />);

    const spinner = screen.getByRole('status');
    expect(spinner).toHaveClass('ds-spinner', 'my-custom-class');
  });

  it('applies a custom label when provided', () => {
    render(<Spinner label="Chargement des données" />);
    const spinner = screen.getByText('Chargement des données');
    expect(spinner).toHaveClass('ds-spinner__label');
  });
});
