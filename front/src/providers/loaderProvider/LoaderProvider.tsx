import { LoaderContext } from "@Front/contexts/loaderContext";
import { Loader } from "@Front/pages/Loader/Loader";
import { type ReactNode, useCallback, useMemo, useState } from "react";

export const LoaderProvider = ({ children }: { children: ReactNode }) => {
  const [isLoading, setIsLoading] = useState(false);

  const showLoader = useCallback(() => setIsLoading(true), []);
  const hideLoader = useCallback(() => setIsLoading(false), []);
  const loaderContextValue = useMemo(
    () => ({ showLoader, hideLoader }),
    [showLoader, hideLoader],
  );

  return (
    <LoaderContext.Provider value={loaderContextValue}>
      {children}
      {isLoading && <Loader />}
    </LoaderContext.Provider>
  );
};
