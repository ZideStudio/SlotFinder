import type { RouteHandle } from '@Front/routing/routeHandle';
import { useMemo } from 'react';
import { Outlet, useMatches, type UIMatch } from 'react-router';
import { Grid } from '../Grid/Grid';
import { Header } from './Header/Header';

export const Layout = () => {
  const matches = useMatches() as UIMatch<unknown, RouteHandle>[];

  const hideHeader = useMemo(
    () => matches.some(match => match.handle?.hideHeader === true),
    [matches],
  );

  return (
    <div>
      {!hideHeader && <Header />}
      <Grid component="main" container>
        <Outlet />
      </Grid>
    </div>
  );
};
