import { render } from '@testing-library/react';
import { ClickIcon } from '../ClickIcon';

describe('ClickIcon', () => {
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

  it('applies props to the button element', () => {
    const onClick = vitest.fn();
    const { getByRole } = render(
      <ClickIcon
        icon={props => (
          <svg {...props}>
            <rect width="100" height="100" fill="blue" />
          </svg>
        )}
        onClick={onClick}
      />,
    );
    const button = getByRole('button');
    button.click();
    expect(onClick).toHaveBeenCalled();
  });
});
