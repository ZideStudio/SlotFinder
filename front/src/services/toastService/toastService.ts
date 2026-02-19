import { AbstractObserver } from "@Front/observers/abstractObserver";

export type Toast = {
  id: number;
  message: string;
  duration?: number;
};

export class ToastService extends AbstractObserver {
  private currentId: number = 0;

  private readonly defaultDuration: number;

  private readonly toast: Map<number, Toast> = new Map();

  private cachedAllToastIds: Array<number> = [];

  constructor(duration: number = 3000) {
    super();

    this.defaultDuration = duration;
  }

  addToast(toast: string, duration?: number) {
    this.currentId += 1;

    const durationToUse = duration || this.defaultDuration;

    this.toast.set(this.currentId, {
      id: this.currentId,
      message: toast,
      duration: durationToUse,
    });

    this.cachedAllToastIds = Array.from(this.toast.keys());

    setTimeout(() => {
      this.removeToast(this.currentId);
    }, durationToUse);

    this.notifyObservers();
  }

  getAllToastIds() {
    return this.cachedAllToastIds;
  }

  getToastById(id: number) {
    return this.toast.get(id);
  }

  removeToast(id: number) {
    this.toast.delete(id);

    this.cachedAllToastIds = Array.from(this.toast.keys());

    this.notifyObservers();
  }
}
