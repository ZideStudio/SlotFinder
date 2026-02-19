export type VoidFunction = () => void;

export abstract class AbstractObserver {
  private readonly observers: Set<VoidFunction> = new Set();

  subscribe(observer: VoidFunction): VoidFunction {
    this.observers.add(observer);

    return () => {
      this.observers.delete(observer);
    };
  }

  protected notifyObservers(): void {
    this.observers.forEach(observer => observer());
  }
}
