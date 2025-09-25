// oxlint-disable no-magic-numbers
import { render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { DebugGrid } from '../DebugGrid';

describe('DebugGrid', () => {
  it('should render the checkbox checked when isCheckedByDefault is true', () => {
    render(<DebugGrid isCheckedByDefault />);
    const cardCheckbox = screen.getByRole('checkbox');
    expect(cardCheckbox).toBeInTheDocument();
    expect(cardCheckbox).toBeChecked();
  });

  it('should render the checkbox unchecked by default', () => {
    render(<DebugGrid />);
    const cardCheckbox = screen.getByRole('checkbox');
    expect(cardCheckbox).toBeInTheDocument();
    expect(cardCheckbox).not.toBeChecked();
  });

  it('should not render DebugGrid when the environment is not development', () => {
    const originalEnv = import.meta.env.DEV;
    import.meta.env.DEV = false;
    render(<DebugGrid />);
    const cardCheckbox = screen.queryByRole('checkbox');
    expect(cardCheckbox).not.toBeInTheDocument();
    import.meta.env.DEV = originalEnv;
  });

  it('should render the correct number of columns when cols prop is set', () => {
    render(<DebugGrid cols={5} />);
    const grid = screen.getByRole('presentation');
    const cols = grid.querySelectorAll('.cols');
    expect(cols).toHaveLength(5);
  });
});
