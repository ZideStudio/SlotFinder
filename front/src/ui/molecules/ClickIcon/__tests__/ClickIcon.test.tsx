import { render, screen } from '@testing-library/react';
import { ClickIcon } from '../ClickIcon';

describe('ClickIcon', () => {
  it('applies the custom class name', () => {
    render(
      <ClickIcon
        icon={props => (
          <svg {...props}>
            <rect width="100" height="100" fill="blue" />
          </svg>
        )}
        className="custom-class"
      />,
    );
    const button = screen.getByRole('button');
    expect(button).toHaveClass('custom-class');
  });

  it('applies props to the button element', () => {
    const onClick = vi.fn();
    render(
      <ClickIcon
        icon={props => (
          <svg {...props}>
            <rect width="100" height="100" fill="blue" />
          </svg>
        )}
        onClick={onClick}
      />,
    );
    const button = screen.getByRole('button');
    button.click();
    expect(onClick).toHaveBeenCalledTimes(1);
  });
});
