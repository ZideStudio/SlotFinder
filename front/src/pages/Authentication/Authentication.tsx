import { Outlet } from 'react-router';
import { OAuth } from './OAuth';

export const Authentication = () => (
  <section>
    <h1>Authentication Page</h1>
    <Outlet />
    <OAuth />
  </section>
);
