import { Grid } from '@Front/components/Grid/Grid';
import { Outlet } from 'react-router';
import { OAuth } from './OAuth';

export const Authentication = () => (
  <Grid component="section" container colSpan={12}>
    <h1>Authentication Page</h1>
    <Grid colSpan={2} colStart={6}>
      <Outlet />
      <OAuth />
    </Grid>
  </Grid>
);
