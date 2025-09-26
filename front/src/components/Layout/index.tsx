import { Outlet } from 'react-router';
import { DebugGrid } from '../DebugGrid/DebugGrid';

export const Layout = () => (
  <>
    <main className="grid">
      <Outlet />
    </main>
    <DebugGrid />
  </>
);
