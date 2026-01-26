import { AuthenticationContext } from '@Front/contexts/AuthenticationContext/AuthenticationContext';
import { useContext } from 'react';

export const useAuthenticationContext = () => {
  const authenticationContext = useContext(AuthenticationContext);

  if (!authenticationContext) {
    throw new Error('useAuthenticationContext must be used within a AuthenticationContextProvider');
  }

  return authenticationContext;
};
