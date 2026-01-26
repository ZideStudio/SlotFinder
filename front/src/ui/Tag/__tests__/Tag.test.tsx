import { render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { Tag } from '../Tag';

describe('Tag', () => {
  it('should render tag component with filled style', () => {
    render(<Tag color="#007bff">Text</Tag>);
    const tag = screen.getByText('Text');
    const parentTag = tag.parentElement;
    expect(tag).toBeInTheDocument();
    expect(tag).toHaveAttribute('title', 'Text');
    expect(parentTag).toHaveClass('ds-tag--filled');
    expect(parentTag).toHaveStyle({ backgroundColor: '#007bff' });
    expect(parentTag).toHaveStyle({ borderColor: '#007bff' });
  });

  it('should render tag component with outlined style', () => {
    render(
      <Tag color="#007bff" appearance="outlined">
        Text
      </Tag>,
    );
    const tag = screen.getByText('Text');
    const parentTag = tag.parentElement;
    expect(tag).toBeInTheDocument();
    expect(tag).toHaveAttribute('title', 'Text');
    expect(parentTag).toHaveClass('ds-tag--outlined');
    expect(parentTag).toHaveStyle({ borderColor: '#007bff' });
  });
});
