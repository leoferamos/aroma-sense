import { createContext } from 'react';
import type { User } from '../types/auth';

export type UserRole = 'admin' | 'client';

export interface AuthContextType {
  user: User | null;
  role: UserRole | null;
  accessToken: string | null;
  isAuthenticated: boolean;
  isReady: boolean;
  setAuth: (accessToken: string, user: User) => void;
  refreshUser: () => Promise<void>;
  logout: () => Promise<void>;
}

export const AuthContext = createContext<AuthContextType | undefined>(undefined);
