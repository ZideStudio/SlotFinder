import { render, screen } from '@testing-library/react';
import { Header } from '../Header';

describe('Header', () => {
  it('renders the header with logo and buttons', () => {
    render(<Header />);

    const logo = screen.getByAltText('Slot Finder logo');
    expect(logo).toBeInTheDocument();

    const buttons = screen.getAllByRole('button');
    expect(buttons).toHaveLength(2);
  });
});
