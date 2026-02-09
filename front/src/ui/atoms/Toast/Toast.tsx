import React from 'react';
import { getClassName } from '@Front/utils/getClassName';
import './Toast.scss';

type ToastProps = {
  children: React.ReactNode;
  onClose?: () => void;
};

export const Toast = ({ children, onClose }: ToastProps) => {
  const parentClassName = getClassName({
    defaultClassName: 'ds-toast',
    modifiers: ['visible'],
  });

  return (
    <div className={parentClassName} role="status" aria-live="polite">
      <span className="ds-toast__content">{children}</span>

      {onClose && (
        <button className="ds-toast__close" onClick={onClose} aria-label="Fermer">
          ✕
        </button>
      )}
    </div>
  );
};
