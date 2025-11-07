import { useState } from 'react';

interface UseForgotPasswordReturn {
  requestPasswordReset: (email: string) => Promise<void>;
  loading: boolean;
  error: string | null;
  success: string | null;
}

export const useForgotPassword = (): UseForgotPasswordReturn => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  const requestPasswordReset = async (email: string) => {
    setLoading(true);
    setError(null);
    setSuccess(null);

    try {
      // endpoint placeholder — replace with real endpoint later
      const res = await fetch('/api/auth/forgot-password', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email }),
      });

      if (!res.ok) {
        const payload = await res.json().catch(() => ({}));
        throw new Error(payload?.message || 'Failed to request password reset.');
      }

      setSuccess('If the email exists in our system, you will receive instructions to reset your password.');
    } catch (err: any) {
      setError(err?.message || 'Unknown error. Please try again later.');
    } finally {
      setLoading(false);
    }
  };

  return {
    requestPasswordReset,
    loading,
    error,
    success,
  };
};