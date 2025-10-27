import React, { createContext, useContext, useState, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import { logoutUser } from '../services/auth';

type UserRole = 'admin' | 'client';

interface AuthContextType {
  role: UserRole | null;
  setRole: (role: UserRole | null) => void;
  logout: () => Promise<void>;
  isAuthenticated: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [role, setRole] = useState<UserRole | null>(null);
  const navigate = useNavigate();

  const logout = useCallback(async () => {
    try {
      // Call route to clear HttpOnly cookie
      await logoutUser();
    } catch (error) {
      console.error('Logout error:', error);
    } finally {
      // Always clear local state and redirect
      setRole(null);
      navigate('/login');
    }
  }, [navigate]);

  const isAuthenticated = role !== null;

  return (
    <AuthContext.Provider value={{ role, setRole, logout, isAuthenticated }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};
