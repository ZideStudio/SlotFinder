import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { ToastService } from '../toastService';

describe('ToastService', () => {
  it('should return empty array on initialization for getAllToastIds', () => {
    const store = new ToastService();

    expect(store.getAllToastIds()).toStrictEqual([]);
  });

  it('should add a toast and increment currentId with default duration', () => {
    const store = new ToastService();

    store.addToast('Hello World');

    const toastIds = store.getAllToastIds();
    const toast = store.getToastById(toastIds[0]);

    expect(toast).toBeDefined();
    expect(toast?.id).toStrictEqual(toastIds[0]);
    expect(toast?.message).toStrictEqual('Hello World');
    expect(toast?.duration).toStrictEqual(3000);
  });

  it('should accept custom duration when adding a toast', () => {
    const store = new ToastService();

    store.addToast('Short', 1500);

    const toastIds = store.getAllToastIds();
    expect(store.getToastById(toastIds[0])?.duration).toStrictEqual(1500);
  });

  it('should use constructor default duration when provided', () => {
    const store = new ToastService(10000);

    store.addToast('Long');

    const toastIds = store.getAllToastIds();
    expect(store.getToastById(toastIds[0])?.duration).toStrictEqual(10000);
  });

  it('should return all toast ids after adding multiple toasts', () => {
    const store = new ToastService();

    store.addToast('A');
    store.addToast('B');
    store.addToast('C');

    const toastIds = store.getAllToastIds();
    expect(toastIds.length).toBe(3);
  });

  it('should remove toast and update ids', () => {
    const store = new ToastService();

    store.addToast('A');
    store.addToast('B');
    store.addToast('C');

    const toastIds = store.getAllToastIds();
    store.removeToast(toastIds[1]);

    expect(store.getToastById(toastIds[1])).toBeUndefined();
    expect(store.getAllToastIds()).toStrictEqual([toastIds[0], toastIds[2]]);
  });

  it('should notify subscribers on add and remove, and allow unsubscribe', () => {
    const store = new ToastService();

    const observer = vi.fn();
    const unsubscribe = store.subscribe(observer);

    store.addToast('One');
    expect(observer).toHaveBeenCalledTimes(1);

    store.addToast('Two');
    expect(observer).toHaveBeenCalledTimes(2);

    unsubscribe();

    observer.mockClear();

    store.addToast('Three');
    expect(observer).not.toHaveBeenCalled();

    const observer2 = vi.fn();
    const unsubscribe2 = store.subscribe(observer2);

    const toastIds = store.getAllToastIds();
    store.removeToast(toastIds[0]);
    expect(observer2).toHaveBeenCalledTimes(1);

    unsubscribe2();
  });

  it('should notify subscribers on automatic removal', () => {
    vi.useFakeTimers();

    const storeAuto = new ToastService();
    const observerAuto = vi.fn();
    const unsubscribeAuto = storeAuto.subscribe(observerAuto);

    storeAuto.addToast('Auto');
    expect(observerAuto).toHaveBeenCalledTimes(1); // Add

    vi.advanceTimersByTime(3000);
    expect(observerAuto).toHaveBeenCalledTimes(2); // Removal after timeout

    // Verify the toast was actually removed
    const toastIds = storeAuto.getAllToastIds();
    expect(storeAuto.getToastById(toastIds[0])).toBeUndefined();
    expect(storeAuto.getAllToastIds()).toStrictEqual([]);

    unsubscribeAuto();
    vi.useRealTimers();
    vi.clearAllTimers();
  });
});
