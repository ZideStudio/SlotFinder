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
