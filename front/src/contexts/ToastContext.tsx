import React from 'react';

export type ToastItem = {
  id: string;
  message: string;
};

type ToastContextType = {
  show: (message: string) => void;
};

export const ToastContext = React.createContext<ToastContextType | null>(null);
