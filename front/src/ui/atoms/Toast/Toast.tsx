import React, { useEffect, useState } from 'react';
import { getClassName } from '@Front/utils/getClassName';
import './Toast.scss';

type ToastProps = {
  children: React.ReactNode;
  className?: string;
  onClose?: () => void;
};

export const Toast = ({ children, className, onClose }: ToastProps) => {
  const [isVisible, setIsVisible] = useState(false);

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
      <span className="ds-toast__content">{children}</span>

      {onClose && (
        <button className="ds-toast__close" onClick={onClose} aria-label="Fermer">
          ✕
        </button>
      )}
    </div>
  );
};
