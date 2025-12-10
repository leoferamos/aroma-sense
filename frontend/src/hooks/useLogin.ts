import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { isAxiosError } from "axios";
import { loginUser } from "../services/auth";
import { useTranslation } from 'react-i18next';
import { useAuth } from "../hooks/useAuth";
import type { LoginRequest } from "../types/auth";

export function useLogin() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const navigate = useNavigate();
  const { setAuth } = useAuth();
  const { t } = useTranslation('common');

  async function login(data: LoginRequest) {
    setLoading(true);
    setError("");
    try {
      const res = await loginUser(data);
      
      // Save access token and user to AuthContext
      setAuth(res.access_token, res.user);
      
      // Redirect based on role
      if (res.user.role === "admin" || res.user.role === "super_admin") {
        navigate("/admin/dashboard");
      } else {
        navigate("/products");
      }
    } catch (err: unknown) {
      if (isAxiosError(err)) {
        const errorMsg = err.response?.data?.error?.toLowerCase() || "";
        if (errorMsg.includes("invalid credentials") || errorMsg.includes("invalid_credentials")) {
          setError(t('auth.invalidCredentials'));
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

  return { login, loading, error };
}