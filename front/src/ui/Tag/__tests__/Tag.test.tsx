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

  it('should render tag component with filled in decision style', () => {
    render(<Tag variant="inDecision">Text</Tag>);
    const tag = screen.getByText('Text').parentElement;
    expect(tag).toBeInTheDocument();
    expect(tag).toHaveClass('ds-tag--inDecision');
    expect(tag).toHaveClass('ds-tag--filled');
  });

  it('should render tag component with outlined in decision style', () => {
    render(
      <Tag variant="inDecision" appearance="outlined">
        Text
      </Tag>,
    );
    const tag = screen.getByText('Text').parentElement;
    expect(tag).toBeInTheDocument();
    expect(tag).toHaveClass('ds-tag--inDecision');
    expect(tag).toHaveClass('ds-tag--outlined');
  });

  it('should render tag component with filled coming soon style', () => {
    render(<Tag variant="comingSoon">Text</Tag>);
    const tag = screen.getByText('Text').parentElement;
    expect(tag).toBeInTheDocument();
    expect(tag).toHaveClass('ds-tag--comingSoon');
    expect(tag).toHaveClass('ds-tag--filled');
  });

  it('should render tag component with outlined coming soon style', () => {
    render(
      <Tag variant="comingSoon" appearance="outlined">
        Text
      </Tag>,
    );
    const tag = screen.getByText('Text').parentElement;
    expect(tag).toBeInTheDocument();
    expect(tag).toHaveClass('ds-tag--comingSoon');
    expect(tag).toHaveClass('ds-tag--outlined');
  });

  it('should render tag component with filled finished style', () => {
    render(<Tag variant="finished">Text</Tag>);
    const tag = screen.getByText('Text').parentElement;
    expect(tag).toBeInTheDocument();
    expect(tag).toHaveClass('ds-tag--finished');
    expect(tag).toHaveClass('ds-tag--filled');
  });

  it('should render tag component with outlined finished style', () => {
    render(
      <Tag variant="finished" appearance="outlined">
        Text
      </Tag>,
    );
    const tag = screen.getByText('Text').parentElement;
    expect(tag).toBeInTheDocument();
    expect(tag).toHaveClass('ds-tag--finished');
    expect(tag).toHaveClass('ds-tag--outlined');
  });
});
