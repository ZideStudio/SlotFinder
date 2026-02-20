import { type PropsWithChildren, useRef } from 'react';
import { ToastService } from '@Front/services/toastService/toastService';
import { ToastContext } from '@Front/contexts/toastContext';
import { ToastContainer } from './ToastContainer';

export type ToastProviderProps = PropsWithChildren & {
  defaultDuration?: number;
};

export const ToastProvider = ({ children, defaultDuration }: ToastProviderProps) => {
  const toastRef = useRef({ toast: new ToastService(defaultDuration) });


  return (
    <ToastContext.Provider value={toastRef.current}>
      <ToastContainer />
      {children}
    </ToastContext.Provider>
  );
};
