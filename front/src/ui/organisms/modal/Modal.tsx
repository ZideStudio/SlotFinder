// oxlint-disable react/jsx-props-no-spreading
import { Heading } from '@Front/ui/atoms/Heading/Heading';
import { Button } from '@Front/ui/molecules/Button/Button';
import { getClassName } from '@Front/utils/getClassName';
import { forwardRef, useId, type DialogHTMLAttributes, type ReactNode } from 'react';

import './Modal.scss';

type ModalProps = {
  title: string;
  children: ReactNode;
  className?: string;
  action?: () => void;
} & DialogHTMLAttributes<HTMLDialogElement>;

export const Modal = forwardRef<HTMLDialogElement, ModalProps>(
  ({ title, className, children, action, ...props }, ref) => {
    const titleId = useId();
    const parentClassName = getClassName({
      defaultClassName: 'ds-modal',
      className,
    });

    const closeModal = () => {
      if (ref && 'current' in ref && ref.current) {
        ref.current.close();
      }
    };

    return (
      <dialog aria-labelledby={titleId} className={parentClassName} ref={ref} {...props} closedby="any">
        <div className="ds-modal__header">
          <Heading level={1} id={titleId}>
            {title}
          </Heading>
          <Button aria-label="Fermer la fenêtre" onClick={closeModal} className="ds-modal__close-button">
            x
          </Button>
        </div>

        {children}

        <div className="ds-modal__footer">
          <Button onClick={closeModal} className="ds-modal__close-button">
            Fermer
          </Button>
          <Button onClick={action} className="ds-modal__close-action">
            Action
          </Button>
        </div>
      </dialog>
    );
  },
);

Modal.displayName = 'Modal';
