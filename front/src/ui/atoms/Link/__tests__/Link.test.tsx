import { Link } from '../Link';
import { render, screen } from '@testing-library/react';

describe('Link', () => {
  it('should render a link with required href and children props', () => {
    render(<Link href="https://example.com">Example Link</Link>);
    const link = screen.getByRole('link', { name: 'Example Link' });
    expect(link).toBeInTheDocument();
    expect(link).toHaveAttribute('href', 'https://example.com');
    expect(link).toHaveClass('ds-link');
  });

  it('should open link in a new tab when openInNewTab is true', () => {
    render(
      <Link href="https://example.com" openInNewTab>
        New Tab Link
      </Link>,
    );
    const link = screen.getByRole('link', { name: 'New Tab Link' });
    expect(link).toHaveAttribute('target', '_blank');
    expect(link).toHaveAttribute('rel', 'noopener noreferrer');
  });

  it('should apply custom className', () => {
    render(
      <Link href="https://example.com" className="custom-class">
        Custom Class Link
      </Link>,
    );
    const link = screen.getByRole('link', { name: 'Custom Class Link' });
    expect(link).toHaveClass('ds-link custom-class');
  });

  it('should not have target and rel attributes when openInNewTab is false', () => {
    render(<Link href="https://example.com">Same Tab Link</Link>);
    const link = screen.getByRole('link', { name: 'Same Tab Link' });
    expect(link).not.toHaveAttribute('target');
    expect(link).not.toHaveAttribute('rel');
  });

  it('should render children correctly', () => {
    render(
      <Link href="https://example.com">
        <span>Child Element</span>
      </Link>,
    );
    const childElement = screen.getByText('Child Element');
    expect(childElement).toBeInTheDocument();
  });
});
