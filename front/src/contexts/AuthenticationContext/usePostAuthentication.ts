import { useState, type Dispatch, type SetStateAction } from 'react';

export type UsePostAuthenticationReturn = {
  postAuthRedirectPath: string | undefined;
  setPostAuthRedirectPath: Dispatch<SetStateAction<string | undefined>>;
  resetPostAuthRedirectPath: () => void;
};

export const usePostAuthentication = (): UsePostAuthenticationReturn => {
  const [postAuthRedirectPath, setPostAuthRedirectPath] = useState<string>();

  const resetPostAuthRedirectPath = () => {
    setPostAuthRedirectPath(undefined);
  };

  return {
    postAuthRedirectPath,
    setPostAuthRedirectPath,
    resetPostAuthRedirectPath,
  };
};
