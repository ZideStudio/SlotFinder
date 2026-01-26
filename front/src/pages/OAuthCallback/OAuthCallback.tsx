import { appRoutes } from '@Front/routing/appRoutes';
import { isInternalUrl } from '@Front/utils/isInternalUrl';
import { Navigate, useSearchParams } from 'react-router';

export const OAuthCallback = () => {
  const [searchParams] = useSearchParams();
  const returnUrl = searchParams.get('returnUrl');
  const error = searchParams.get('error');

  if (error) {
    return <Navigate to={appRoutes.error()} state={{ message: error }} replace />;
  }

  const destinationPath = returnUrl && isInternalUrl(returnUrl) ? returnUrl : appRoutes.home();

  return <Navigate to={destinationPath} replace />;
};
