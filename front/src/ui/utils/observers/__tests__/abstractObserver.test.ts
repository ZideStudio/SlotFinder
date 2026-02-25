import { describe, expect, it, vi } from 'vitest';
import { AbstractObserver } from '../abstractObserver';

class AbstractObserverTest extends AbstractObserver {
  public triggerNotify(): void {
    this.notifyObservers();
  }
}

describe('AbstractObserver', () => {
  it('should add a listener and call it on notifyListeners', () => {
    const observerTest = new AbstractObserverTest();
    const listener = vi.fn();
    observerTest.subscribe(listener);
    observerTest.triggerNotify();

    expect(listener).toHaveBeenCalledTimes(1);
  });

  it('should remove listener when unsubscribe is called', () => {
    const observerTest = new AbstractObserverTest();
    const listener = vi.fn();
    const unsubscribe = observerTest.subscribe(listener);
    unsubscribe();

    observerTest.triggerNotify();

    expect(listener).not.toHaveBeenCalled();
  });

  it('should only notify remaining listeners after unsubscribe is called', () => {
    const observerTest = new AbstractObserverTest();
    const listener1 = vi.fn();
    const listener2 = vi.fn();
    const listener3 = vi.fn();
    const unsubscribe = observerTest.subscribe(listener1);
    observerTest.subscribe(listener2);
    observerTest.subscribe(listener3);
    unsubscribe();

    observerTest.triggerNotify();

    expect(listener1).not.toHaveBeenCalled();
    expect(listener2).toHaveBeenCalledTimes(1);
    expect(listener3).toHaveBeenCalledTimes(1);
  });
});
