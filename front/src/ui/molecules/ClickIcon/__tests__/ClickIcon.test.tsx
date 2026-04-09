import { render } from '@testing-library/react';
import { ClickIcon } from '../ClickIcon';

describe('ClickIcon', () => {
  it('renders the icon correctly', () => {
    render(
      <ClickIcon
        icon={props => (
          <svg {...props}>
            <rect width="100" height="100" fill="blue" />
          </svg>
        )}
      />,
    );
  });

  it('applies the custom class name', () => {
    const { container } = render(
      <ClickIcon
        icon={props => (
          <svg {...props}>
            <rect width="100" height="100" fill="blue" />
          </svg>
        )}
        className="custom-class"
      />,
    );
    expect(container.firstChild).toHaveClass('custom-class');
  });
});
