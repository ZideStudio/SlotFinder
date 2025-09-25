import '@Front/i18n';
import { StrictMode } from 'react';
import { createRoot, type Root } from 'react-dom/client';
import { App } from './components/App/App';

// eslint-disable-next-line import/no-default-export
export default class Bootstrap extends HTMLElement {
  private readonly root: Root;

  private isMounted = false;

  constructor() {
    super();

    this.root = createRoot(this);
  }

  connectedCallback() {
    this.isMounted = true;
    this.root.render(
      <StrictMode>
        <App basename={this.getAttribute('basename') ?? ''} />
      </StrictMode>,
    );
  }

  disconnectedCallback() {
    this.isMounted = false;

    queueMicrotask(() => {
      if (!this.isMounted) {
        this.root.unmount();
      }
    });
  }
}
