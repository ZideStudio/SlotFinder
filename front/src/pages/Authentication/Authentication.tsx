import { Grid } from '@Front/components/Grid/Grid';
import { Outlet } from 'react-router';
import { OAuth } from './OAuth';

export const Authentication = () => (
  <Grid
    component="section"
    container
    colSpan={{
      'desktop-small': 12,
      tablet: 8,
      mobile: 4,
    }}
  >
    <Grid
      component="h1"
      colSpan={{
        'desktop-small': 12,
        tablet: 8,
        mobile: 4,
      }}
    >
      Authentication Page
    </Grid>
    <Grid
      colSpan={{ 'desktop-small': 2, tablet: 2, mobile: 4 }}
      colStart={{ 'desktop-small': 6, tablet: 4, mobile: 1 }}
    >
      <Outlet />
    </Grid>
    <OAuth />
  </Grid>
);
