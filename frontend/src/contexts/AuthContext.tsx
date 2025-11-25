import React, { createContext, useState, useCallback, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { logoutUser, refreshToken } from '../services/auth';
import { setAccessToken } from '../services/api';
import type { User } from '../types/auth';

type UserRole = 'admin' | 'client';

interface AuthContextType {
  user: User | null;
  role: UserRole | null;
  accessToken: string | null;
  isAuthenticated: boolean;
  isReady: boolean; // Indicates if initial session check is complete
  setAuth: (accessToken: string, user: User) => void;
  logout: () => Promise<void>;
}
const AuthContext = createContext<AuthContextType | undefined>(undefined);

export { AuthContext };
export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null);
  const [accessToken, setAccessTokenState] = useState<string | null>(null);
  const [isReady, setIsReady] = useState(false);
  const navigate = useNavigate();

  const updateAccessToken = useCallback((token: string | null) => {
    setAccessTokenState(token);
    setAccessToken(token);
  }, []);
  const setAuth = useCallback((token: string, userData: User) => {
    updateAccessToken(token);
    setUser(userData);
  }, [updateAccessToken]);


  const logout = useCallback(async () => {
    try {
      await logoutUser();
    } catch (error) {
      console.error('Logout error:', error);
    } finally {
      updateAccessToken(null);
      setUser(null);
      navigate('/login');
    }
  }, [navigate, updateAccessToken]);

  useEffect(() => {
    let isMounted = true;

    (async () => {
      try {
        const data = await refreshToken();
        if (isMounted) {
          updateAccessToken(data.access_token);
          setUser(data.user);
        }
      } catch (error) {
        console.error('Token refresh failed:', error);
        // No valid refresh token cookie
        if (isMounted) {
          updateAccessToken(null);
          setUser(null);
        }
      } finally {
        if (isMounted) {
          setIsReady(true);
        }
      }
    })();

    return () => {
      isMounted = false;
    };
  }, [updateAccessToken]);

  const role = user?.role as UserRole | null;
  const isAuthenticated = user !== null && accessToken !== null;

  return (
    <AuthContext.Provider value={{
      user,
      role,
      accessToken,
      isAuthenticated,
      isReady,
      setAuth,
      logout
    }}>
      {children}
    </AuthContext.Provider>
  );
};
