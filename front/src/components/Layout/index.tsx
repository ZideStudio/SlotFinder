import { Outlet } from "react-router";
import { Header } from "./Header/Header";

export const Layout = () => (
  <>
    <Header />
    <main className="grid">
      <Outlet />
    </main>
  </>
);
