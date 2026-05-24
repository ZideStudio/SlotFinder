import { Outlet } from 'react-router';
import { Grid } from '../Grid/Grid';
import { Header } from './Header/Header';

export const Layout = () => {
  return (
    <div>
      <Header />
      <Grid component="main" container>
        <Outlet />
      </Grid>
    </div>
  );
};
