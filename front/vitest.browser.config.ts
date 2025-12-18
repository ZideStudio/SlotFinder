import react from '@vitejs/plugin-react-swc';
import { playwright } from '@vitest/browser-playwright';
import { resolve } from 'path';
import { loadEnv } from 'vite';
import svgr from 'vite-plugin-svgr';
import viteTsconfigPaths from 'vite-tsconfig-paths';
// eslint-disable-next-line import/no-unresolved
import { defineConfig } from 'vitest/config';

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '');
  return {
    plugins: [
      react(),
      viteTsconfigPaths({
        projects: [resolve(__dirname, './tsconfig.test.json')],
      }),
      svgr(),
    ],
    envPrefix: env.ENV_PREFIX ?? 'FRONT_',
    test: {
      globals: true,
      // TODO: setup vitest browser with msw : https://mswjs.io/docs/recipes/vitest-browser-mode
      clearMocks: true,
      css: false,
      reporters: ['default', 'junit', 'vitest-sonar-reporter'],
      outputFile: {
        'vitest-sonar-reporter': 'sonar-report.xml',
        junit: 'junit-report.xml',
      },
      include: ['src/**/*.browser.test.[jt]sx'],
      poolOptions: {
        forks: {
          minForks: env.CI ? 1 : undefined,
          maxForks: env.CI ? 2 : undefined,
        },
      },
      coverage: {
        enabled: env.CI === 'true',
        reporter: ['lcovonly', 'html', 'text', 'text-summary'],
        provider: 'v8',
        lines: 80,
        functions: 75,
        branches: 80,
        statements: 80,
        include: ['src/**/*.[jt]s?(x)'],
        exclude: [
          'src/**/*.d.[jt]s?(x)',
          'src/utils/tests/**/*.[jt]s?(x)',
          'src/**/*.types.[jt]s?(x)',
          'src/**/*.stories.[jt]s?(x)',
          'src/**/*.(spec|test|steps).[jt]s?(x)',
          'src/**/__tests__/**/*.[jt]s?(x)',
          'src/i18n/**',
          'src/routing/**',
          'src/main.ts',
        ],
      },
      browser: {
        enabled: true,
        provider: playwright(),
        // https://vitest.dev/config/browser/playwright
        instances: [{ browser: 'chromium' }],
        screenshotFailures: false,
      },
    },
  };
});
