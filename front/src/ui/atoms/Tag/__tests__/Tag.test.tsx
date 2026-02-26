import { render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { Tag } from '../Tag';

describe('Tag', () => {
  it('should render tag component with filled style', () => {
    render(<Tag color="#007bff">Text</Tag>);
    const tag = screen.getByText('Text');
    expect(tag).toBeInTheDocument();
    expect(tag).toHaveAttribute('title', 'Text');
    expect(tag).toHaveClass('ds-tag--filled');
    expect(tag).toHaveStyle({ '--tag-color': '#007bff' });
  });

  it('should render tag component with outlined style', () => {
    render(
      <Tag color="#007bff" appearance="outlined">
        Text
      </Tag>,
    );
    const tag = screen.getByText('Text');
    expect(tag).toBeInTheDocument();
    expect(tag).toHaveAttribute('title', 'Text');
    expect(tag).toHaveClass('ds-tag--outlined');
    expect(tag).toHaveStyle({ '--tag-color': '#007bff' });
  });

  it('should render tag component with custom class name', () => {
    render(
      <Tag color="#007bff" className="custom-class">
        Text
      </Tag>,
    );
    const tag = screen.getByText('Text');
    expect(tag).toBeInTheDocument();
    expect(tag).toHaveClass('custom-class');
  });

  it('should render tag component with title attribute', () => {
    render(
      <Tag color="#007bff" title="Custom title">
        Text
      </Tag>,
    );
    const tag = screen.getByText('Text');
    expect(tag).toBeInTheDocument();
    expect(tag).toHaveAttribute('title', 'Custom title');
  });

  it('should render tag component with children text as title when title prop is not provided', () => {
    render(<Tag color="#007bff">Text</Tag>);
    const tag = screen.getByText('Text');
    expect(tag).toBeInTheDocument();
    expect(tag).toHaveAttribute('title', 'Text');
  });
});
