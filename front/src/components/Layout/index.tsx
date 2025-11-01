import { Outlet } from 'react-router';
import { Grid } from '../Grid/Grid';

export const Layout = () => (
  <Grid component="main" container>
    <Outlet />
  </Grid>
);
