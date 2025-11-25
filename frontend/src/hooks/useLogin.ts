import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { isAxiosError } from "axios";
import { loginUser } from "../services/auth";
import { messages } from "../constants/messages";
import { useAuth } from "../hooks/useAuth";
import type { LoginRequest } from "../types/auth";

export function useLogin() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const navigate = useNavigate();
  const { setAuth } = useAuth();

  async function login(data: LoginRequest) {
    setLoading(true);
    setError("");
    try {
      const res = await loginUser(data);
      
      // Save access token and user to AuthContext
      setAuth(res.access_token, res.user);
      
      // Redirect based on role
      if (res.user.role === "admin") {
        navigate("/admin/dashboard");
      } else {
        navigate("/products");
      }
    } catch (err: unknown) {
      if (isAxiosError(err)) {
        const errorMsg = err.response?.data?.error?.toLowerCase() || "";
        if (errorMsg.includes("invalid credentials")) {
          setError(messages.invalidCredentials || "Invalid email or password.");
        } else {
          setError(err.response?.data?.error || messages.genericError);
        }
      } else if (err instanceof Error) {
        setError(err.message || messages.genericError);
      } else {
        setError(messages.genericError);
      }
    } finally {
      setLoading(false);
    }
  }

  return { login, loading, error };
}