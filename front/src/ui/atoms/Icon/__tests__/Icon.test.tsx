import { render } from '@testing-library/react';
import { Icon } from '../Icon';

describe('Icon', () => {
  it('should render the icon component', () => {
    const { container } = render(
      <Icon
        icon={() => (
          <svg>
            <rect width="100" height="100" fill="blue" />
          </svg>
        )}
      />,
    );

    expect(container.querySelector('.ds-icon')).toBeInTheDocument();
    expect(container.querySelector('svg')).toBeInTheDocument();
  });

  it('should apply custom class name', () => {
    const { container } = render(
      <Icon
        icon={() => (
          <svg>
            <rect width="100" height="100" fill="blue" />
          </svg>
        )}
        className="custom-class"
      />,
    );

    expect(container.querySelector('.ds-icon')).toHaveClass('custom-class');
  });
});
