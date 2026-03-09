import { render, screen } from '@testing-library/react';
import { Icon } from '../Icon';

describe('Icon', () => {
  it('should render the icon component', () => {
    render(
      <Icon
        icon={props => (
          <svg {...props}>
            <rect width="100" height="100" fill="blue" />
          </svg>
        )}
      />,
    );

    expect(screen.getByRole('presentation', { hidden: true })).toBeInTheDocument();
  });

  it('should apply custom class name', () => {
    render(
      <Icon
        icon={props => (
          <svg {...props}>
            <rect width="100" height="100" fill="blue" />
          </svg>
        )}
        className="custom-class"
      />,
    );

    expect(screen.getByRole('presentation', { hidden: true })).toHaveClass('custom-class');
  });
});
