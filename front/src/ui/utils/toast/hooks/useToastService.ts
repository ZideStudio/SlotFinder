import { ToastContext } from '@Front/contexts/toastContext';
import { useContext } from 'react';

export const useToastService = () => {
  const context = useContext(ToastContext);

  if (!context) {
    throw new Error('useToast must be used within a ToastProvider');
  }

  return context.toast;
};
