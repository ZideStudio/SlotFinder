import { AuthenticationContext } from "@Front/contexts/AuthenticationContext/AuthenticationContext";
import type { AuthenticationContextType } from "@Front/contexts/AuthenticationContext/types";
import { useContext } from "react";

export const useAuthenticationContext = () => useContext(AuthenticationContext) as AuthenticationContextType;