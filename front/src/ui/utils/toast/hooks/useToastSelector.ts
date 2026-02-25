import { useSyncExternalStore } from 'react';
import { useToastService } from './useToastService';
import type { ToastService } from '@Front/ui/utils/toast/service/toastService/toastService';

export const useToastSelector = <T>(selector: (toastService: ToastService) => T) => {
  const toastService = useToastService();

  return useSyncExternalStore(
    listener => toastService.subscribe(listener),
    () => selector(toastService),
    () => selector(toastService),
  );
};
