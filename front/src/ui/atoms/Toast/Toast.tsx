import { useToastSelector } from '@Front/ui/utils/toast/hooks/useToastSelector';
import { useToastService } from '@Front/ui/utils/toast/hooks/useToastService';
import { getClassName } from '@Front/utils/getClassName';
import { memo } from 'react';

import './Toast.scss';

type ToastProps = {
  className?: string;
  id: string;
};

const ToastComponent = ({ className, id }: ToastProps) => {
  const toast = useToastSelector(t => t.getToastById(id));
  const toastService = useToastService();

  if (!toast) {
    return null;
  }

  const parentClassName = getClassName({
    defaultClassName: 'ds-toast',
    className,
  });

  return (
    <div className={parentClassName} role="status">
      <span className="ds-toast__content">{toast.message}</span>
      <button
        className="ds-toast__close"
        onClick={() => {
          toastService.removeToast(id);
        }}
        aria-label="Fermer la notification"
      >
        ✕
      </button>
    </div>
  );
};

export const Toast = memo(ToastComponent);
