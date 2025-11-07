import { appRoutes } from '@Front/routing/appRoutes';
import { isInternalUrl } from '@Front/utils/isInternalUrl';
import { Navigate, useSearchParams } from 'react-router';

export const OAuthCallback = () => {
  const [searchParams] = useSearchParams();
  const returnUrl = searchParams.get('returnUrl');
  const message = searchParams.get('message');

  if (message) {
    return <Navigate to={appRoutes.error()} state={{ message }} replace />;
  }

  const destinationPath = returnUrl && isInternalUrl(returnUrl) ? returnUrl : appRoutes.dashboard();

  return <Navigate to={destinationPath} replace />;
};
