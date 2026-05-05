import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import type { SVGProps } from 'react';
import { ClickIcon } from '../ClickIcon';

const TestIcon = (props: SVGProps<SVGSVGElement>) => (
  <svg {...props}>
    <rect width="100" height="100" fill="blue" />
  </svg>
);

describe('ClickIcon', () => {
  it('applies the custom class name', () => {
    render(<ClickIcon icon={TestIcon} className="custom-class" />);
    const button = screen.getByRole('button');
    expect(button).toHaveClass('custom-class');
  });

  it('applies props to the button element', async () => {
    const onClick = vi.fn();
    render(<ClickIcon icon={TestIcon} onClick={onClick} />);
    const button = screen.getByRole('button');
    await userEvent.click(button);
    expect(onClick).toHaveBeenCalledTimes(1);
  });

  it('renders as anchor when as="a"', () => {
    render(<ClickIcon as="a" icon={TestIcon} href="/test" />);
    const anchor = screen.getByRole('link');
    expect(anchor).toHaveAttribute('href', '/test');
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
