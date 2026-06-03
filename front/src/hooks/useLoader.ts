import { LoaderContext } from "@Front/contexts/loaderContext";
import { useContext } from "react";

export const useLoader = () => {
  const context = useContext(LoaderContext);

  if (!context) {
    throw new Error("useLoader must be used within a LoaderProvider");
  }

  const { showLoader, hideLoader } = context;

  return { showLoader, hideLoader };
};
