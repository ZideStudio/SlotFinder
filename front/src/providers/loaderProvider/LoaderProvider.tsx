import { LoaderContext } from "@Front/contexts/loaderContext";
import { Loader } from "@Front/pages/Loader/Loader";
import React from "react";

export const LoaderProvider = ({ children }: { children: React.ReactNode }) => {
  const [isLoading, setIsLoading] = React.useState(false);

  const showLoader = React.useCallback(() => setIsLoading(true), []);
  const hideLoader = React.useCallback(() => setIsLoading(false), []);
  const loaderContextValue = React.useMemo(
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
