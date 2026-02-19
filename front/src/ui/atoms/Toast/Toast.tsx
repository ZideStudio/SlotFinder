import React, { memo, useEffect, useState } from 'react';
import { getClassName } from '@Front/utils/getClassName';
import './Toast.scss';
import { useToastSelector } from '@Front/hooks/useToastSelector';
import { useToastService } from '@Front/hooks/useToastService';

type ToastProps = {
  className?: string;
  id: number;
};

const ToastComponent = ({ className, id }: ToastProps) => {
  const [isVisible, setIsVisible] = useState(false);
  const toast = useToastSelector(t => t.getToastById(id));
  const toastService = useToastService();

  if (!toast) {
    return null;
  }

  useEffect(() => {
    requestAnimationFrame(() => {
      setIsVisible(true);
    });
  }, []);

  const parentClassName = getClassName({
    defaultClassName: 'ds-toast',
    modifiers: [isVisible ? 'visible' : '', className || ''],
  });

  return (
    <div className={parentClassName} role="status" aria-live="polite">
      <span className="ds-toast__content">{toast.message}</span>
      <button
        className="ds-toast__close"
        onClick={() => {
          toastService.removeToast(id);
        }}
        aria-label="Fermer"
      >
        ✕
      </button>
    </div>
  );
};

export const Toast = memo(ToastComponent);
