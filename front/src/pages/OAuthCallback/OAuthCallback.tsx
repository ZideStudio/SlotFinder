import { Navigate, useSearchParams } from 'react-router-dom';

export const OAuthCallback = () => {
  const [searchParams] = useSearchParams();
  const message = searchParams.get('message');

  if (message) {
    return <Navigate to="/error" state={{ message }} replace />;
  }

  return <Navigate to="/" replace />;
};
