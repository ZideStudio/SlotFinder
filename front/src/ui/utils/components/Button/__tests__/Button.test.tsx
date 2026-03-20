import { render, screen } from '@testing-library/react';
import { Button } from '../Button';

const CustomButton = ({ children }: { children: React.ReactNode }) => <span>{children}</span>;

describe('Button', () => {
  it('renders children', () => {
    render(<Button>Button</Button>);
    expect(screen.getByText('Button')).toBeInTheDocument();
  });

  it('renders as anchor when as="a"', () => {
    render(
      <Button as="a" href="/test">
        Button
      </Button>,
    );
    const anchorElement = screen.getByText('Button').closest('a');
    expect(anchorElement).not.toBeNull();
    expect(anchorElement).toHaveAttribute('href', '/test');
  });

  it('renders as custom component when as is a component', () => {
    render(<Button as={CustomButton}>Button</Button>);
    const customElement = screen.getByText('Button').closest('span');
    expect(customElement).not.toBeNull();
  });

  it('applies variant and color classes', () => {
    render(
      <Button variant="secondary" color="danger">
        Button
      </Button>,
    );
    const buttonElement = screen.getByText('Button');
    expect(buttonElement).toHaveClass('ds-button--secondary');
    expect(buttonElement).toHaveClass('ds-button--danger');
  });

  it('applies disabled class when disabled', () => {
    render(<Button disabled>Button</Button>);
    const buttonElement = screen.getByText('Button');
    expect(buttonElement).toHaveClass('ds-button--disabled');
  });

  it('renders icon when icon prop is provided', () => {
    const TestIcon = () => <svg data-testid="test-icon" />;
    render(<Button icon={TestIcon}>Button</Button>);
    expect(screen.getByTestId('test-icon')).toBeInTheDocument();
  });

  it('applies additional class names', () => {
    render(<Button className="custom-class">Button</Button>);
    expect(screen.getByText('Button')).toHaveClass('custom-class');
  });
});
