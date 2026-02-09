import { ToastContext } from '@Front/contexts/ToastContext';
import { useContext } from 'react';

export const useToastContext = () => {
  const toastContext = useContext(ToastContext);

  if (!toastContext) {
    throw new Error(
      'useToastContext must be used within ToastProvider'
    );
  }

  return toastContext;
};
