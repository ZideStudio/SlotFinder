import { defineConfig } from 'vitest/config';
import { getBaseConfig } from './base';

/**
 * Workspace config that runs unit and browser tests in a single Vitest process.
 * Coverage results from both projects are automatically merged into one report.
 *
 * Usage:
 *   npm test               → all tests (unit + browser)
 *   npm test -- --coverage → all tests with merged coverage report
 */
export default defineConfig(({ mode }) => {
  const {
    test: { coverage },
  } = getBaseConfig(mode);

  return {
    test: {
      coverage,
      projects: ['./config/vitest/vitest.unit.config.ts', './config/vitest/vitest.browser.config.ts'],
    },
  };
});
