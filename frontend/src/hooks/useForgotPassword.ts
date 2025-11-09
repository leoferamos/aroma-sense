import { useState } from 'react';
import { requestPasswordReset, confirmPasswordReset } from '../services/auth';

interface UseForgotPasswordReturn {
  requestReset: (email: string) => Promise<boolean>;
  confirmReset: (email: string, code: string, newPassword: string) => Promise<boolean>;
  loading: boolean;
  error: string | null;
  success: string | null;
}

export const useForgotPassword = (): UseForgotPasswordReturn => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  const requestReset = async (email: string): Promise<boolean> => {
    setLoading(true);
    setError(null);
    setSuccess(null);

    try {
      const response = await requestPasswordReset(email);
      setSuccess(response.message);
      return true;
    } catch (err: any) {
      setError(err?.response?.data?.error || 'Failed to request password reset. Please try again.');
      return false;
    } finally {
      setLoading(false);
    }
  };

  const confirmReset = async (email: string, code: string, newPassword: string): Promise<boolean> => {
    setLoading(true);
    setError(null);
    setSuccess(null);

    try {
      const response = await confirmPasswordReset(email, code, newPassword);
      setSuccess(response.message);
      return true;
    } catch (err: any) {
      setError(err?.response?.data?.error || 'Invalid or expired code. Please try again.');
      return false;
    } finally {
      setLoading(false);
    }
  };

  return {
    requestReset,
    confirmReset,
    loading,
    error,
    success,
  };
};