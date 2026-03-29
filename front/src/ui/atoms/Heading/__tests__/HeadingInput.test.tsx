// oxlint-disable no-magic-numbers
import { render } from '@testing-library/react';
import { Heading } from '../Heading';

describe('Heading', () => {
  it.each([1, 2, 3] as const)('renders the correct heading level %i', (level: 1 | 2 | 3) => {
    const { getByText } = render(<Heading level={level}>Test Heading</Heading>);
    const headingElement = getByText('Test Heading');
    expect(headingElement.tagName).toBe(`H${level}`);
  });

  it('applies the correct class names', () => {
    const { getByText } = render(
      <Heading level={1} className="custom-class">
        Test Heading
      </Heading>,
    );
    const headingElement = getByText('Test Heading');
    expect(headingElement).toHaveClass('ds-heading');
    expect(headingElement).toHaveClass('ds-heading--level-1');
    expect(headingElement).toHaveClass('custom-class');
  });
});
