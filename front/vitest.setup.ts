/* eslint-disable import/no-extraneous-dependencies */
import { server } from '@Mocks/server';
import '@testing-library/jest-dom/vitest';
import { cleanup } from '@testing-library/react';
import * as matchers from 'vitest-axe/matchers';

expect.extend(matchers);

beforeAll(() => {
  server.listen({ onUnhandledRequest: 'error' });

  vi.mock('react-i18next', () => ({
    useTranslation: vi.fn().mockReturnValue({
      t: (messageId: string, args: Record<string, unknown>) => messageId + (args ? `::${JSON.stringify(args)}` : ''),
    }),
    initReactI18next: {
      type: '3rdParty',
      init: () => {},
    },
  }));
});

afterAll(() => {
  server.close();
});

afterEach(() => {
  server.resetHandlers();
  cleanup();
});
