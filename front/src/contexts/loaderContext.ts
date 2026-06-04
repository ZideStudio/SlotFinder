import { createContext } from "react";

export type LoaderContextValue = {
  showLoader: () => void;
  hideLoader: () => void;
};

export const LoaderContext = createContext<LoaderContextValue | null>(null);
