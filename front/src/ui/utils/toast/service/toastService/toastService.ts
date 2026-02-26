import { AbstractObserver } from "@Front/ui/utils/observers/abstractObserver";

export type Toast = {
  id: string;
  message: string;
  duration?: number;
};

export class ToastService extends AbstractObserver {
  private readonly defaultDuration: number;

  private readonly toast = new Map<string, Toast>();

  private cachedAllToastIds: string[] = [];

  constructor(duration: number = 3000) {
    super();

    this.defaultDuration = duration;
  }

  addToast(toast: string, duration?: number) {
    const newId = globalThis.crypto.randomUUID();

    const durationToUse = duration || this.defaultDuration;

    this.toast.set(newId, {
      id: newId,
      message: toast,
      duration: durationToUse,
    });

    this.cachedAllToastIds = Array.from(this.toast.keys());

    setTimeout(() => {
      this.removeToast(newId);
    }, durationToUse);

    this.notifyObservers();
  }

  getAllToastIds() {
    return this.cachedAllToastIds;
  }

  getToastById(id: string) {
    return this.toast.get(id);
  }

  removeToast(id: string) {
    this.toast.delete(id);

    this.cachedAllToastIds = Array.from(this.toast.keys());

    this.notifyObservers();
  }
}