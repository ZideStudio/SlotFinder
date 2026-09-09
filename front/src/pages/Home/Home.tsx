import { useAuthenticationContext } from "@Front/hooks/useAuthenticationContext";
import { Dashboard } from "./Dashboard/Dashboard";
import { Welcome } from "./Welcome/Welcome";

export const Home = () => {
  const { isAuthenticated } = useAuthenticationContext();
  return isAuthenticated ? <Dashboard /> : <Welcome />;
};
