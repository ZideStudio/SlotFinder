import { render } from '@testing-library/react';
import { Heading } from '../Heading';

describe('Heading', () => {
  it('renders the correct heading level 1', () => {
    const { getByText } = render(<Heading level={1}>Test Heading</Heading>);
    const headingElement = getByText('Test Heading');
    expect(headingElement.tagName).toBe('H1');
  });

  it('renders the correct heading level 2', () => {
    const { getByText } = render(<Heading level={2}>Test Heading</Heading>);
    const headingElement = getByText('Test Heading');
    expect(headingElement.tagName).toBe('H2');
  });

  it('renders the correct heading level 3', () => {
    const { getByText } = render(<Heading level={3}>Test Heading</Heading>);
    const headingElement = getByText('Test Heading');
    expect(headingElement.tagName).toBe('H3');
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
