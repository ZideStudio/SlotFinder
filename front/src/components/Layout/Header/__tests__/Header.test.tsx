import { routeObject } from '@Front/routing/routes';
import { renderRoute } from '@Front/utils/testsUtils/customRender/customRender';
import { getAuthStatus200 } from '@Mocks/handlers/authStatusHandlers';
import { server } from '@Mocks/server';
import { render, screen } from '@testing-library/react';
import { Header } from '../Header';

describe('Header', () => {
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
  it('renders the header with logo and buttons', () => {
    renderRoute({
      initialEntry: '/',
      routes: [
        {
          path: '/',
          element: <Header />,
        },
      ],
    });

    const logo = screen.getByAltText('Slot Finder logo');
    expect(logo).toBeInTheDocument();

    const buttons = screen.getAllByRole('button');
    expect(buttons).toHaveLength(2);
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
