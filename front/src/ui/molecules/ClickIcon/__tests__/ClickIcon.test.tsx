import { render, screen } from '@testing-library/react';
import { ClickIcon } from '../ClickIcon';

const TestIcon = (props: React.SVGProps<SVGSVGElement>) => (
  <svg {...props}>
    <rect width="100" height="100" fill="blue" />
  </svg>
);

const CustomSpan = ({ children, ...props }: { children: React.ReactNode }) => (
  <span data-testid="custom-span" {...props}>{children}</span>
);

describe('ClickIcon', () => {
  it('applies the custom class name', () => {
    render(<ClickIcon icon={TestIcon} className="custom-class" />);
    const button = screen.getByRole('button');
    expect(button).toHaveClass('custom-class');
  });

  it('applies props to the button element', () => {
    const onClick = vi.fn();
    render(<ClickIcon icon={TestIcon} onClick={onClick} />);
    const button = screen.getByRole('button');
    button.click();
    expect(onClick).toHaveBeenCalledTimes(1);
  });

  it('renders as anchor when as="a"', () => {
    render(<ClickIcon as="a" icon={TestIcon} href="/test" />);
    const anchor = screen.getByRole('link');
    expect(anchor).toHaveAttribute('href', '/test');
  });

  it('renders as custom component when as is a component', () => {
    render(<ClickIcon as={CustomSpan} icon={TestIcon} />);
    expect(screen.getByTestId('custom-span')).toBeInTheDocument();
  });

  describe('type attribute', () => {
    it('has type="button" by default', () => {
      render(<ClickIcon icon={TestIcon} />);
      expect(screen.getByRole('button')).toHaveAttribute('type', 'button');
    });

    it('has type="button" when as="button"', () => {
      render(<ClickIcon as="button" icon={TestIcon} />);
      expect(screen.getByRole('button')).toHaveAttribute('type', 'button');
    });

    it('respects an explicit type prop', () => {
      render(<ClickIcon type="submit" icon={TestIcon} />);
      expect(screen.getByRole('button')).toHaveAttribute('type', 'submit');
    });

    it('does not set type when rendered as an anchor', () => {
      render(<ClickIcon as="a" icon={TestIcon} href="/test" />);
      expect(screen.getByRole('link')).not.toHaveAttribute('type');
    });
  });
});
