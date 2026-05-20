import { renderRoute } from '@Front/utils/testsUtils/customRender/customRender';
import { getAuthStatus200 } from '@Mocks/handlers/authStatusHandlers';
import { server } from '@Mocks/server';
import { screen } from '@testing-library/react';
import { routeObject } from '../../../routing/routes';

describe('Layout', () => {
  type TestRoute = '/withHeader' | '/withoutHeader';

  const renderLayoutRoute = (initialEntry: TestRoute) =>
    renderRoute({
      initialEntry,
      routes: [
        {
          ...routeObject[0],
          index: false,
          children: [
            {
              path: '/withHeader',
              element: <p>with header</p>,
            },
            {
              path: '/withoutHeader',
              element: <p>without header</p>,
              handle: {
                hideHeader: true,
              },
            },
          ],
        },
      ],
    });

  describe('when the route has no hideHeader handle', () => {
    beforeEach(() => {
      server.use(getAuthStatus200);
    });
    it('should render the header banner', async () => {
      renderLayoutRoute('/withHeader');

      await expect(screen.findByRole('banner')).resolves.toBeInTheDocument();
    });
  });

  describe('when the route has hideHeader: true', () => {
    it('should not render the header banner', async () => {
      renderLayoutRoute('/withoutHeader');

      await screen.findByText('without header');

      expect(screen.queryByRole('banner')).not.toBeInTheDocument();
    });
  });
});
