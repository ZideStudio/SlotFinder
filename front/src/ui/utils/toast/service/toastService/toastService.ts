import { AbstractObserver } from "@Front/ui/utils/observers/abstractObserver";

export type Toast = {
  id: string;
  message: string;
  duration?: number;
  timeout?: ReturnType<typeof setTimeout>;
};

export class ToastService extends AbstractObserver {
  private readonly defaultDuration: number | null;

  private readonly toast = new Map<string, Toast>();

  private cachedAllToastIds: string[] = [];

  constructor(duration: number | null = 3000) {
    super();

    this.defaultDuration = duration;
  }

  addToast(toast: string, duration?: number | null) {
    const newId = globalThis.crypto.randomUUID();

    const durationToUse = duration === undefined ? this.defaultDuration : duration;

    let timeout: ReturnType<typeof setTimeout> | undefined;
    if (durationToUse !== null) {
      timeout = setTimeout(() => {
        this.removeToast(newId);
      }, durationToUse);
    }

    this.toast.set(newId, {
      id: newId,
      message: toast,
      duration: durationToUse ?? undefined,
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