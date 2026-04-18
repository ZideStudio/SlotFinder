import { renderHook } from '@testing-library/react';
import { type RefObject, createRef } from 'react';
import { useModal } from '../useModal';

describe('useModal', () => {
  describe('when no existing ref is provided', () => {
    it('should initialize modalRef with a null current value', () => {
      const { result } = renderHook(() => useModal());

      expect(result.current.modalRef.current).toBeNull();
    });
  });

  describe('when an existing ref is provided', () => {
    it('should use the provided ref instead of the default one', () => {
      const existingRef: RefObject<HTMLDialogElement | null> = createRef();
      const { result } = renderHook(() => useModal(existingRef));

      expect(result.current.modalRef).toStrictEqual(existingRef);
    });
  });

  describe('openModal', () => {
    it('should call showModal on the dialog element', () => {
      const showModal = vi.fn();
      const mockDialog = { showModal, close: vi.fn() } as unknown as HTMLDialogElement;
      const existingRef = { current: mockDialog } as RefObject<HTMLDialogElement>;

      const { result } = renderHook(() => useModal(existingRef));

      result.current.openModal();

      expect(showModal).toHaveBeenCalledTimes(1);
    });

    it('should do nothing when modalRef.current is null', () => {
      const existingRef = { current: null } as RefObject<HTMLDialogElement | null>;
      const { result } = renderHook(() => useModal(existingRef));

      expect(() => result.current.openModal()).not.toThrow();
    });
  });

  describe('closeModal', () => {
    it('should call close on the dialog element', () => {
      const close = vi.fn();
      const mockDialog = { showModal: vi.fn(), close } as unknown as HTMLDialogElement;
      const existingRef = { current: mockDialog } as RefObject<HTMLDialogElement>;
      const { result } = renderHook(() => useModal(existingRef));

      result.current.closeModal();

      expect(close).toHaveBeenCalledTimes(1);
    });

    it('should do nothing when modalRef.current is null', () => {
      const existingRef = { current: null } as RefObject<HTMLDialogElement | null>;
      const { result } = renderHook(() => useModal(existingRef));

      expect(() => result.current.closeModal()).not.toThrow();
    });
  });
});
