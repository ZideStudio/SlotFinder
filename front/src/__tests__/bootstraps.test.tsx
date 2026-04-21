import { App } from '@Front/components/App/App';
import { waitFor } from '@testing-library/react';
import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { type Mock } from 'vitest';
import bootstrap from '../bootstrap';

// oxlint-disable-next-line vitest/prefer-import-in-mock jest/prefer-ending-with-an-expect
vi.mock('@Front/components/App', () => ({
  App: () => <div data-testid="app-mock">AppMock</div>,
}));

// oxlint-disable-next-line vitest/prefer-import-in-mock jest/prefer-ending-with-an-expect
vi.mock('react-dom/client', () => ({
  createRoot: vi.fn(),
}));

describe('bootstrap', () => {
  // oxlint-disable-next-line init-declarations
  let container: HTMLElement;
  const render = vi.fn();
  const unmount = vi.fn();

  beforeEach(() => {
    render.mockReset();
    unmount.mockReset();
    (createRoot as Mock).mockReset();
    (createRoot as Mock).mockReturnValue({
      render,
      unmount,
    });
    if (!customElements.get('bootstrap-html-element')) {
      customElements.define('bootstrap-html-element', bootstrap);
    }
    container = document.createElement('bootstrap-html-element');
  });

  afterEach(() => {
    container.remove();
  });

  it('should mount the React component in the custom element', () => {
    document.body.appendChild(container);

    expect(createRoot as Mock).toHaveBeenCalledWith(expect.anything());
    expect(render).toHaveBeenCalledWith(
      <StrictMode>
        <App basename="" />
      </StrictMode>,
    );
    expect(unmount).not.toHaveBeenCalled();
  });

  it('should unmount the React component when removed from the DOM', async () => {
    document.body.appendChild(container);

    expect(unmount).not.toHaveBeenCalled();

    container.remove();

    await waitFor(() => {
      expect(unmount).toHaveBeenCalledWith();
    });
  });

  it('should not unmount when element is removed and re-added to the DOM', () => {
    document.body.appendChild(container);
    container.remove();
    document.body.appendChild(container);

    expect(unmount).not.toHaveBeenCalled();
  });
});
