import { AUTH_STORAGE_KEYS } from '../constants/auth';
import type { User } from '../types/auth';

// Auth utility functions
export const getStoredUser = (): User | null => {
  try {
    const userStr = localStorage.getItem(AUTH_STORAGE_KEYS.USER);
    return userStr ? JSON.parse(userStr) : null;
  } catch (error) {
    console.error('Error parsing stored user:', error);
    return null;
  }
};

export const setStoredUser = (user: User): void => {
  try {
    localStorage.setItem(AUTH_STORAGE_KEYS.USER, JSON.stringify(user));
  } catch (error) {
    console.error('Error storing user:', error);
  }
};

export const getStoredToken = (): string | null => {
  return localStorage.getItem(AUTH_STORAGE_KEYS.ACCESS_TOKEN);
};

export const setStoredToken = (token: string): void => {
  localStorage.setItem(AUTH_STORAGE_KEYS.ACCESS_TOKEN, token);
};

export const getStoredRefreshToken = (): string | null => {
  return localStorage.getItem(AUTH_STORAGE_KEYS.REFRESH_TOKEN);
};

export const setStoredRefreshToken = (token: string): void => {
  localStorage.setItem(AUTH_STORAGE_KEYS.REFRESH_TOKEN, token);
};

export const clearStoredAuth = (): void => {
  localStorage.removeItem(AUTH_STORAGE_KEYS.USER);
  localStorage.removeItem(AUTH_STORAGE_KEYS.ACCESS_TOKEN);
  localStorage.removeItem(AUTH_STORAGE_KEYS.REFRESH_TOKEN);
};

export const isTokenExpired = (token: string): boolean => {
  try {
    const payload = JSON.parse(atob(token.split('.')[1]));
    const currentTime = Date.now() / 1000;
    return payload.exp < currentTime;
  } catch (error) {
    console.error('Error checking token expiration:', error);
    return true;
  }
};