import { errorRoutes } from '@Front/pages/Error';
import { appRoutes } from '@Front/routing/appRoutes';
import { renderRoute } from '@Front/utils/testsUtils/customRender';
import { screen } from '@testing-library/react';
import type { InitialEntry } from 'react-router-dom';
import { describe, expect, it } from 'vitest';

const renderErrorPage = (message?: string) => {
  const initialEntries: InitialEntry[] = [];

  if (message) {
    initialEntries.push({ pathname: appRoutes.error(), state: { message } });
  } else {
    initialEntries.push(appRoutes.error());
  }

  return renderRoute({
    routes: [errorRoutes],
    routesOptions: { initialEntries },
  });
};

describe('ErrorPage', () => {
  it('should display the error message when provided in state', async () => {
    renderErrorPage('TestError');
    expect(await screen.findByRole('alert')).toHaveTextContent('TestError');
  });

  it('should display the default unexpected message when no message is provided', async () => {
    renderErrorPage();
    expect(await screen.findByRole('alert')).toHaveTextContent('unexpected');
  });
});
