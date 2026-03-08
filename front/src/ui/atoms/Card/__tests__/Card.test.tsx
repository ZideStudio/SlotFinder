import { render, screen } from '@testing-library/react';
import { Card } from '../Card';

describe('Card', () => {
  it('renders children', () => {
    render(<Card>Card</Card>);
    expect(screen.getByText('Card')).toBeInTheDocument();
  });

  it('applies additional class names', () => {
    render(<Card className="custom-class">Card</Card>);
    expect(screen.getByText('Card')).toHaveClass('custom-class');
  });
});
