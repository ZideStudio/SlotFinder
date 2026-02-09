import React from 'react';
import { getClassName } from '@Front/utils/getClassName';

import './Toast.scss';

export type ToastProps = {
  className?: string;
  children: React.ReactNode;
  visible?: boolean;
  onClose?: () => void;
};

export const Toast = ({ children, className, visible = true, onClose }: ToastProps) => {
  if (!visible) return null;

  const parentClassName = getClassName({
    defaultClassName: 'ds-toast',
    className,
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
