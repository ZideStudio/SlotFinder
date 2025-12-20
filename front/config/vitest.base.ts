import react from '@vitejs/plugin-react-swc';
import { resolve } from 'node:path';
import { loadEnv } from 'vite';
import svgr from 'vite-plugin-svgr';
import viteTsconfigPaths from 'vite-tsconfig-paths';


const MAX_WORKERS = 2;
export const getBaseConfig = (mode: string) => {
  const env = loadEnv(mode, process.cwd(), '');
  return {
    plugins: [
      react(),
      viteTsconfigPaths({
        projects: [resolve(__dirname, '../tsconfig.test.json')],
      }),
      svgr(),
    ],
    envPrefix: env.ENV_PREFIX ?? 'FRONT_',
    test: {
      globals: true,
      clearMocks: true,
      css: false,
      reporters: ['default', 'junit', 'vitest-sonar-reporter'],
      outputFile: {
        'vitest-sonar-reporter': 'sonar-report.xml',
        junit: 'junit-report.xml',
      },
      maxWorkers: env.CI ? MAX_WORKERS : undefined,
      coverage: {
        enabled: env.CI === 'true',
        reporter: ['lcovonly', 'html', 'text', 'text-summary'],
        provider: 'v8' as const,
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
    },
  };
};
