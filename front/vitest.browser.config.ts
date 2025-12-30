import { playwright } from '@vitest/browser-playwright';
import { defineConfig } from 'vitest/config';
import { getBaseConfig } from './config/vitest.base';

// TODO: setup vitest browser with msw : https://mswjs.io/docs/recipes/vitest-browser-mode
export default defineConfig(({ mode }) => {
  const base = getBaseConfig(mode);
  return {
    ...base,
    test: {
      ...base.test,
      include: ['src/**/*.browser.test.[jt]sx'],
      setupFiles: ['./vitest.browser.setup.ts'],
      browser: {
        enabled: true,
        provider: playwright(),
        instances: [{ browser: 'chromium' }],
      },
    },
  };
});
