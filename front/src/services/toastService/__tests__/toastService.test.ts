import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { ToastService } from '../toastService';

beforeEach(() => {
  vi.useFakeTimers();
});

afterEach(() => {
  vi.useRealTimers();
  vi.clearAllTimers();
});

describe('ToastService', () => {
  it('should return empty array on initialization for getAllToastIds', () => {
    const store = new ToastService();

    expect(store.getAllToastIds()).toStrictEqual([]);
  });

  it('should add a toast and increment currentId with default duration', () => {
    const store = new ToastService();

    store.addToast('Hello World');

    const toast = store.getToastById(1);

    expect(toast).toBeDefined();
    expect(toast?.id).toStrictEqual(1);
    expect(toast?.message).toStrictEqual('Hello World');
    expect(toast?.duration).toStrictEqual(3000);
  });

  it('should accept custom duration when adding a toast', () => {
    const store = new ToastService();

    store.addToast('Short', 1500);

    expect(store.getToastById(1)?.duration).toStrictEqual(1500);
  });

  it('should use constructor default duration when provided', () => {
    const store = new ToastService(10000);

    store.addToast('Long');

    expect(store.getToastById(1)?.duration).toStrictEqual(10000);
  });

  it('should return all toast ids after adding multiple toasts', () => {
    const store = new ToastService();

    store.addToast('A');
    store.addToast('B');
    store.addToast('C');

    expect(store.getAllToastIds()).toStrictEqual([1, 2, 3]);
  });

  it('should remove toast and update ids', () => {
    const store = new ToastService();

    store.addToast('A');
    store.addToast('B');
    store.addToast('C');

    store.removeToast(2);

    expect(store.getToastById(2)).toBeUndefined();
    expect(store.getAllToastIds()).toStrictEqual([1, 3]);
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

    store.addToast('Three');
    expect(observer).toHaveBeenCalledTimes(2);

    const observer2 = vi.fn();
    const unsubscribe2 = store.subscribe(observer2);

    store.removeToast(1);
    expect(observer2).toHaveBeenCalledTimes(1);

    unsubscribe2();
  });

  it('should notify subscribers on automatic removal', () => {
    const storeAuto = new ToastService();
    const observerAuto = vi.fn();
    const unsubscribeAuto = storeAuto.subscribe(observerAuto);

    storeAuto.addToast('Auto');
    expect(observerAuto).toHaveBeenCalledTimes(1); // add

    vi.advanceTimersByTime(3000);
    expect(observerAuto).toHaveBeenCalledTimes(2); // removal after timeout

    // verify the toast was actually removed
    expect(storeAuto.getToastById(1)).toBeUndefined();
    expect(storeAuto.getAllToastIds()).toStrictEqual([]);

    unsubscribeAuto();
  });
});
