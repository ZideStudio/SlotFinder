import { useRef, type RefObject } from 'react';

type UseModalReturnType = {
  modalRef: RefObject<HTMLDialogElement | null>;
  openModal: () => void;
  closeModal: () => void;
};

export const useModal = (existingModalRef?: RefObject<HTMLDialogElement | null>): UseModalReturnType => {
  const defaultModalRef = useRef<HTMLDialogElement>(null);
  const modalRef = existingModalRef ?? defaultModalRef;

  const openModal = () => {
    if (!modalRef.current) {
      return;
    }

    modalRef.current.showModal();
  };

  const closeModal = () => {
    if (!modalRef.current) {
      return;
    }

    modalRef.current.close();
  };

  return { modalRef, openModal, closeModal };
};
