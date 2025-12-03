import React, { useState, useCallback, useEffect } from 'react';
import { AuthContext } from './AuthContextData';
import type { UserRole } from './AuthContextData';
import { useNavigate } from 'react-router-dom';
import { logoutUser, refreshToken } from '../services/auth';
import { setAccessToken } from '../services/api';
import type { User } from '../types/auth';

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

  // Refresh user profile
  const refreshUser = useCallback(async () => {
    try {
      const data = await refreshToken();
      updateAccessToken(data.access_token);
      setUser(data.user);
    } catch {
      updateAccessToken(null);
      setUser(null);
    }
  }, [updateAccessToken]);

  // Listen for account deletion events emitted by API interceptor
  useEffect(() => {
    type AccountEventDetail = {
      deletion_requested_at?: string | null;
      deletion_confirmed_at?: string | null;
      deactivated_at?: string | null;
      deactivated_by?: string | null;
      deactivation_reason?: string | null;
      deactivation_notes?: string | null;
      suspension_until?: string | null;
      contestation_deadline?: string | null;
    };

    const onRequested = (e: Event) => {
      const detail = (e as CustomEvent<AccountEventDetail>).detail;
      setUser((prev) => {
        if (!prev) return prev;
        return { ...prev, deletion_requested_at: detail.deletion_requested_at } as User;
      });
    };
    const onConfirmed = (e: Event) => {
      const detail = (e as CustomEvent<AccountEventDetail>).detail;
      setUser((prev) => {
        if (!prev) return prev;
        return { ...prev, deletion_confirmed_at: detail.deletion_confirmed_at } as User;
      });
    };
    const onDeactivated = (e: Event) => {
      const detail = (e as CustomEvent<AccountEventDetail>).detail;
      setUser((prev) => {
        if (!prev) return prev;
        return {
          ...prev,
          deactivated_at: detail.deactivated_at ?? prev.deactivated_at,
          deactivated_by: detail.deactivated_by ?? prev.deactivated_by,
          deactivation_reason: detail.deactivation_reason ?? prev.deactivation_reason,
          deactivation_notes: detail.deactivation_notes ?? prev.deactivation_notes,
          suspension_until: detail.suspension_until ?? prev.suspension_until,
          contestation_deadline: detail.contestation_deadline ?? prev.contestation_deadline,
        } as User;
      });
    };

    window.addEventListener('account-deletion-requested', onRequested);
    window.addEventListener('account-deletion-confirmed', onConfirmed);
    window.addEventListener('account-deactivated', onDeactivated);
    return () => {
      window.removeEventListener('account-deletion-requested', onRequested);
      window.removeEventListener('account-deletion-confirmed', onConfirmed);
      window.removeEventListener('account-deactivated', onDeactivated);
    };
  }, [navigate]);

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
      refreshUser,
      logout
    }}>
      {children}
    </AuthContext.Provider>
  );
};
