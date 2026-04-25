// oxlint-disable react/jsx-props-no-spreading
import { Heading } from '@Front/ui/atoms/Heading/Heading';
import { Button } from '@Front/ui/molecules/Button/Button';
import { ClickIcon } from '@Front/ui/molecules/ClickIcon/ClickIcon';
import { getClassName } from '@Front/utils/getClassName';
import Close from '@material-symbols/svg-400/rounded/close.svg?react';
import { useId, type ComponentProps, type ComponentPropsWithoutRef, type ReactNode, type RefObject } from 'react';

import { useModal } from '@Front/ui/utils/hooks/useModal';
import './Modal.scss';

type ButtonProps = ComponentProps<typeof Button>;

type ModalProps = {
  title: string;
  children: ReactNode;
  primaryButtonProps: ButtonProps;
  secondaryButtonProps?: ButtonProps;
  ref?: RefObject<HTMLDialogElement | null>;
} & Omit<ComponentPropsWithoutRef<'dialog'>, 'children' | 'title'>;

export const Modal = ({
  title,
  ref,
  className,
  children,
  primaryButtonProps,
  secondaryButtonProps,
  ...props
}: ModalProps) => {
  const { closeModal, modalRef } = useModal(ref);

  const titleId = useId();
  const parentClassName = getClassName({
    defaultClassName: 'ds-modal',
    className,
  });

  return (
    <dialog aria-labelledby={titleId} className={parentClassName} ref={modalRef} {...props} closedby="any">
      <header className="ds-modal__header">
        <Heading level={1} id={titleId}>
          {title}
        </Heading>
        <ClickIcon
          aria-label="Fermer la fenêtre"
          onClick={closeModal}
          className="ds-modal__button--close"
          icon={Close}
          type="button"
        />
      </header>
      
      <main>{children}</main>

      <footer className="ds-modal__footer">
        {secondaryButtonProps && <Button type="button" {...secondaryButtonProps} />}
        <Button type="button" {...primaryButtonProps} />
      </footer>
    </dialog>
  );
};

Modal.displayName = 'Modal';
