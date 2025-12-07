import { useState } from 'react';
import { requestPasswordReset, confirmPasswordReset } from '../services/auth';
import { useTranslation } from 'react-i18next';

interface ApiError {
  response?: {
    data?: {
      error?: string;
    };
  };
}

interface UseForgotPasswordReturn {
  requestReset: (email: string) => Promise<boolean>;
  confirmReset: (email: string, code: string, newPassword: string) => Promise<boolean>;
  loading: boolean;
  error: string | null;
  emailSentSuccess: string | null;
  passwordResetSuccess: string | null;
}

export const useForgotPassword = (): UseForgotPasswordReturn => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [emailSentSuccess, setEmailSentSuccess] = useState<string | null>(null);
  const [passwordResetSuccess, setPasswordResetSuccess] = useState<string | null>(null);
  const { t } = useTranslation('common');

  const requestReset = async (email: string): Promise<boolean> => {
    setLoading(true);
    setError(null);
    setEmailSentSuccess(null);
    setPasswordResetSuccess(null);

    try {
      const response = await requestPasswordReset(email);
      setEmailSentSuccess(response.message);
      return true;
    } catch (err: unknown) {
      const apiError = err as ApiError;
      const errorMessage = apiError?.response?.data?.error;
      
      // Handle specific error messages
      if (errorMessage === 'Too many reset requests. Please try again later.') {
        setError(t('errors.tooManyResetRequests'));
      } else {
        setError(errorMessage || 'Failed to request password reset. Please try again.');
      }
      
      return false;
    } finally {
      setLoading(false);
    }
  };

  const confirmReset = async (email: string, code: string, newPassword: string): Promise<boolean> => {
    setLoading(true);
    setError(null);
    setEmailSentSuccess(null);
    setPasswordResetSuccess(null);

    try {
      const response = await confirmPasswordReset(email, code, newPassword);
      setPasswordResetSuccess(response.message);
      return true;
    } catch (err: unknown) {
      const apiError = err as ApiError;
      setError(apiError?.response?.data?.error || t('errors.invalidOrExpiredCode'));
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
    emailSentSuccess,
    passwordResetSuccess,
  };
};