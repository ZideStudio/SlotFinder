# ⚙️ Config

The `config` folder centralizes shared tool configurations for the project. Each sub-folder is dedicated to a specific tool and contains everything needed to run it — base settings, project-specific overrides, and test environment setups.

## 📑 Table of Contents

- [Folder Organization](#folder-organization)
- [Vitest](#vitest)
- [Learn More](#learn-more)

## <span id="folder-organization">Folder Organization</span>

```
config/
  vitest/
    base.ts                    — Shared Vitest + Vite base configuration
    workspace.ts               — Workspace entry point (merges unit + browser projects)
    vitest.unit.config.ts      — Unit test project configuration (jsdom)
    vitest.browser.config.ts   — Browser test project configuration (Playwright/Chromium)
    setup.unit.ts              — Global setup for unit tests (MSW server, jest-dom, i18n mock)
    setup.browser.ts           — Global setup for browser tests (MSW service worker, i18n)
```

## <span id="vitest">Vitest</span>

The `vitest/` folder organises all Vitest configuration into a single place. It uses the [Vitest workspace / projects](https://vitest.dev/guide/projects) feature to run unit tests and browser tests in a single process, with coverage automatically merged across both.

### How it works

```
workspace.ts
├── vitest.unit.config.ts   (environment: jsdom)
└── vitest.browser.config.ts  (environment: Playwright/Chromium via @vitest/browser)
```

`workspace.ts` is the root config. It delegates test execution to the two project configs and owns the shared coverage settings. When `npm test` is run, Vitest collects coverage from both environments and produces a single merged report.

Coverage is enabled only when the `CI` environment variable is set to `"true"` (standard in GitHub Actions, GitLab CI, etc.).

### Available npm scripts

| Script                       | Description                     |
| ---------------------------- | ------------------------------- |
| `npm test`                   | Run all tests (unit + browser)  |
| `npm run test:unit`          | Run unit tests only             |
| `npm run test:unit:watch`    | Run unit tests in watch mode    |
| `npm run test:browser`       | Run browser tests only          |
| `npm run test:browser:watch` | Run browser tests in watch mode |

### `base.ts`

Exports `getBaseConfig(mode)`, which returns the shared Vite plugins and Vitest options inherited by both project configs:

- **Plugins**: `@vitejs/plugin-react-swc`, `vite-tsconfig-paths`, `vite-plugin-svgr`
- **Reporters**: `default`, `junit`, `vitest-sonar-reporter`
- **Coverage**: v8 provider, `lcov`/`html`/`text` reporters, thresholds at 80% lines/branches/statements and 75% functions

### `setup.unit.ts`

Global setup file executed before each unit test suite. It:

- Starts the MSW Node.js server (`server.listen`)
- Extends `expect` with `jest-dom` and `vitest-axe` matchers
- Mocks `react-i18next` globally so translation keys are returned as-is
- Resets handlers and cleans up the DOM after each test

### `setup.browser.ts`

Global setup file executed before each browser test suite. It:

- Starts the MSW service worker (`worker.start`) before all tests
- Loads global CSS tokens and the i18n configuration
- Resets handlers after each test and stops the worker after all tests

## <span id="learn-more">Learn More</span>

- [Vitest Documentation](https://vitest.dev)
- [Vitest Projects](https://vitest.dev/guide/projects)
- [Vitest Browser Mode](https://vitest.dev/guide/browser)
- [MSW — Vitest Browser Mode recipe](https://mswjs.io/docs/recipes/vitest-browser-mode)
- [@vitest/coverage-v8](https://vitest.dev/guide/coverage)
