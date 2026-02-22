import { useToastSelector } from '@Front/hooks/useToastSelector';
import { Toast } from '@Front/ui/atoms/Toast/Toast';

import './ToastContainer.scss';

export const ToastContainer = () => {
  const toastIds = useToastSelector(toast => toast.getAllToastIds());

  if (toastIds.length === 0) {
    return null;
  }

  return (
    <section className="ds-toast-container">
      {toastIds.map(id => (
        <Toast id={id} key={id} />
      ))}
    </section>
  );
};
