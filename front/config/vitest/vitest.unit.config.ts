import { defineConfig } from 'vitest/config';
import { getBaseConfig } from './base';

export default defineConfig(({ mode }) => {
  const base = getBaseConfig(mode);
  return {
    ...base,
    test: {
      ...base.test,
      name: 'unit',
      root: new URL('../../', import.meta.url).pathname,
      environment: 'jsdom',
      setupFiles: 'config/vitest/setup.unit.ts',
      include: ['src/**/*.(spec|test|steps).[jt]s?(x)'],
      exclude: ['src/**/*.browser.test.[jt]sx'],
    },
  };
});
