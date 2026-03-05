import { render } from '@testing-library/react';
import { Card } from '../Card';

describe('Card', () => {
  it('renders children', () => {
    const { getByText } = render(<Card>Card</Card>);
    expect(getByText('Card')).toBeInTheDocument();
  });

  it('applies additional class names', () => {
    const { container } = render(<Card className="custom-class">Card</Card>);
    expect(container.firstChild).toHaveClass('custom-class');
  });
});
