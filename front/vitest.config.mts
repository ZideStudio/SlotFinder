import { defineConfig } from 'vitest/config';
import { getBaseConfig } from './config/vitest.base';

export default defineConfig(({ mode }) => {
  const base = getBaseConfig(mode);
  return {
    ...base,
    test: {
      ...base.test,
      environment: 'jsdom',
      setupFiles: 'vitest.setup.ts',
      include: ['src/**/*.(spec|test|steps).[jt]s?(x)'],
      exclude: ['src/**/*.browser.test.[jt]sx'],
    },
  };
});
