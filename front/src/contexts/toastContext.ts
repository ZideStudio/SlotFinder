import type { ToastService } from '@Front/services/toastService/toastService';
import { createContext } from 'react';

type ToastContextProps = {
  toast: ToastService;
};

export const ToastContext = createContext<ToastContextProps | null>(null);
