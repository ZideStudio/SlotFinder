import { useModal } from '@Front/ui/utils/hooks/useModal';
import { getClassName } from '@Front/utils/getClassName';
import { useId, type ComponentProps, type ComponentPropsWithoutRef, type ReactNode, type RefObject } from 'react';
import { OverlayContent } from '../OverlayContent/OverlayContent';

import './Modal.scss';

type OverlayContentProps = ComponentProps<typeof OverlayContent>;

type ModalProps = {
  title: string;
  children: ReactNode;
  primaryButtonProps: OverlayContentProps['primaryButtonProps'];
  secondaryButtonProps?: OverlayContentProps['secondaryButtonProps'];
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
      <OverlayContent
        title={title}
        titleId={titleId}
        primaryButtonProps={primaryButtonProps}
        secondaryButtonProps={secondaryButtonProps}
        closeOverlay={closeModal}
      >
        {children}
      </OverlayContent>
    </dialog>
  );
};

Modal.displayName = 'Modal';
