import { createContext } from 'react';
import type { AuthenticationContextType } from './types';

export const AuthenticationContext = createContext<AuthenticationContextType | null>(null);
