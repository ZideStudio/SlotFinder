import React from 'react';
import { Toast } from '@Front/ui/atoms/Toast/Toast';
import { ToastContext, type ToastItem } from '@Front/contexts/ToastContext';

import './ToastProvider.scss';

const TOAST_DURATION = 3000;

export const ToastProvider = ({ children }: { children: React.ReactNode }) => {
  const [toasts, setToasts] = React.useState<ToastItem[]>([]);

  const show = (message: string) => {
    const id = crypto.randomUUID();

    setToasts(prevToasts => [...prevToasts, { id, message }]);

    setTimeout(() => {
      setToasts(prevToasts => prevToasts.filter(toast => toast.id !== id));
    }, TOAST_DURATION);
  };

  const remove = (id: string) => {
    setToasts(prevToasts => prevToasts.filter(toast => toast.id !== id));
  };

  return (
    <ToastContext.Provider value={{ show }}>
      {children}

      <div className="ds-toast-container">
        {toasts.map(toast => (
          <Toast key={toast.id} onClose={() => remove(toast.id)}>
            {toast.message}
          </Toast>
        ))}
      </div>
    </ToastContext.Provider>
  );
};
