import { AbstractObserver } from "@Front/ui/utils/observers/abstractObserver";

export type Toast = {
  id: string;
  message: string;
  duration?: number;
  timeout: ReturnType<typeof setTimeout>;
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

    const timeout = setTimeout(() => {
      this.removeToast(newId);
    }, durationToUse);

    this.toast.set(newId, {
      id: newId,
      message: toast,
      duration: durationToUse,
      timeout,
    });

    this.cachedAllToastIds = Array.from(this.toast.keys());

    this.notifyObservers();
  }

  getAllToastIds() {
    return this.cachedAllToastIds;
  }

  getToastById(id: string) {
    return this.toast.get(id);
  }

  removeToast(id: string) {
    const toast = this.toast.get(id);
    if (toast) {
      clearTimeout(toast.timeout);
    }

    this.toast.delete(id);

    this.cachedAllToastIds = Array.from(this.toast.keys());

    this.notifyObservers();
  }
}