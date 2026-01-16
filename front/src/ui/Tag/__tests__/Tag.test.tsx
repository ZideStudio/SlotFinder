import { render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { Tag } from '../Tag';

describe('Tag', () => {
  it('should render tag component with filled default style', () => {
    render(<Tag>Text</Tag>);
    const tag = screen.getByText('Text');
    const parentTag = tag.parentElement;
    expect(tag).toBeInTheDocument();
    expect(tag).toHaveAttribute('title', 'Text');
    expect(parentTag).toHaveClass('ds-tag--default');
    expect(parentTag).toHaveClass('ds-tag--filled');
  });

  it('should render tag component with outlined default style', () => {
    render(<Tag appearance="outlined">Text</Tag>);
    const tag = screen.getByText('Text');
    const parentTag = tag.parentElement;
    expect(tag).toBeInTheDocument();
    expect(tag).toHaveAttribute('title', 'Text');
    expect(parentTag).toHaveClass('ds-tag--default');
    expect(parentTag).toHaveClass('ds-tag--outlined');
  });

  it('should render tag component with filled success style', () => {
    render(<Tag variant="success">Text</Tag>);
    const tag = screen.getByText('Text').parentElement;
    expect(tag).toBeInTheDocument();
    expect(tag).toHaveClass('ds-tag--success');
    expect(tag).toHaveClass('ds-tag--filled');
  });

  it('should render tag component with outlined success style', () => {
    render(
      <Tag variant="success" appearance="outlined">
        Text
      </Tag>,
    );
    const tag = screen.getByText('Text').parentElement;
    expect(tag).toBeInTheDocument();
    expect(tag).toHaveClass('ds-tag--success');
    expect(tag).toHaveClass('ds-tag--outlined');
  });

  it('should render tag component with filled warning style', () => {
    render(<Tag variant="warning">Text</Tag>);
    const tag = screen.getByText('Text').parentElement;
    expect(tag).toBeInTheDocument();
    expect(tag).toHaveClass('ds-tag--warning');
    expect(tag).toHaveClass('ds-tag--filled');
  });

  it('should render tag component with outlined warning style', () => {
    render(
      <Tag variant="warning" appearance="outlined">
        Text
      </Tag>,
    );
    const tag = screen.getByText('Text').parentElement;
    expect(tag).toBeInTheDocument();
    expect(tag).toHaveClass('ds-tag--warning');
    expect(tag).toHaveClass('ds-tag--outlined');
  });

  it('should render tag component with filled error style', () => {
    render(<Tag variant="error">Text</Tag>);
    const tag = screen.getByText('Text').parentElement;
    expect(tag).toBeInTheDocument();
    expect(tag).toHaveClass('ds-tag--error');
    expect(tag).toHaveClass('ds-tag--filled');
  });

  it('should render tag component with outlined error style', () => {
    render(
      <Tag variant="error" appearance="outlined">
        Text
      </Tag>,
    );
    const tag = screen.getByText('Text').parentElement;
    expect(tag).toBeInTheDocument();
    expect(tag).toHaveClass('ds-tag--error');
    expect(tag).toHaveClass('ds-tag--outlined');
  });
});
