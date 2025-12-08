import { useState } from "react";
import { isAxiosError } from "axios";
import { registerUser } from "../services/auth";
import { useTranslation } from 'react-i18next';
import type { RegisterRequest, RegisterResponse } from "../types/auth";

export function useRegister() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");
  const { t } = useTranslation('common');

  async function register(data: RegisterRequest) {
    setLoading(true);
    setError("");
    setSuccess("");
    try {
      const res: RegisterResponse = await registerUser(data);
      setSuccess(res.message || t('auth.registerSuccess'));
    } catch (err: unknown) {
      if (isAxiosError(err)) {
        const errorMsg = err.response?.data?.error?.toLowerCase() || "";
        if (errorMsg.includes("email already registered")) {
          setError(t('auth.emailAlreadyRegistered'));
        } else {
          setError(err.response?.data?.error || t('errors.failedToLoadProfile'));
        }
      } else if (err instanceof Error) {
        setError(err.message || t('errors.failedToLoadProfile'));
      } else {
        setError(t('errors.failedToLoadProfile'));
      }
    } finally {
      setLoading(false);
    }
  }

  return { register, loading, error, success };
}
