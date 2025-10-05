import { appRoutes } from '@Front/routing/appRoutes';
import { Navigate, useSearchParams } from 'react-router-dom';

export const OAuthCallback = () => {
  const [searchParams] = useSearchParams();
  const message = searchParams.get('message');

  if (message) {
    return <Navigate to={appRoutes.error()} state={{ message }} replace />;
  }

  return <Navigate to={appRoutes.dashboard()} replace />;
};
