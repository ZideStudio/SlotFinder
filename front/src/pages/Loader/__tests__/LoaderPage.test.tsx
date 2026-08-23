import { render, screen, act } from '@testing-library/react';
import LoaderPage from '../LoaderPage';

describe('LoaderPage', () => {
  beforeEach(() => {
    vi.useFakeTimers();
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  it('does not show the loader immediately on mount', () => {
    render(<LoaderPage />);

    expect(screen.queryByRole('status', { name: 'loading' })).not.toBeInTheDocument();
  });

  it('does not show the loader just before 100ms', () => {
    render(<LoaderPage />);

    act(() => {
      vi.advanceTimersByTime(99);
    });

    expect(screen.queryByRole('status', { name: 'loading' })).not.toBeInTheDocument();
  });

  it('shows the loader after 100ms', () => {
    render(<LoaderPage />);

    act(() => {
      vi.advanceTimersByTime(100);
    });

    expect(screen.getByRole('status', { name: 'loading' })).toBeInTheDocument();
  });
});
