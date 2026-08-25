import { ToastContext } from "@Front/contexts/toastContext";
import { ToastService } from "@Front/ui/utils/toast/service/toastService/toastService";
import { type PropsWithChildren, useMemo } from "react";
import { ToastContainer } from "./ToastContainer";

export type ToastProviderProps = PropsWithChildren & {
  defaultDuration?: number | null;
};

export const ToastProvider = ({
  children,
  defaultDuration,
}: ToastProviderProps) => {
  const toastRef = useMemo(
    () => ({ toast: new ToastService(defaultDuration) }),
    [defaultDuration],
  );

  return (
    <ToastContext.Provider value={toastRef}>
      <ToastContainer />
      {children}
    </ToastContext.Provider>
  );
};
