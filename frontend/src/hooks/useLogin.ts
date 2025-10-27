import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { loginUser } from "../services/auth";
import { messages } from "../constants/messages";
import { useAuth } from "../contexts/AuthContext";
import type { LoginRequest, LoginResponse } from "../types/auth";

export function useLogin() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const navigate = useNavigate();
  const { setRole } = useAuth();

  async function login(data: LoginRequest) {
    setLoading(true);
    setError("");
    try {
      const res: LoginResponse = await loginUser(data);
      
      // Save only role to auth context
      setRole(res.user.role as 'admin' | 'client');
      
      // Redirect based on role
      if (res.user.role === "admin") {
        navigate("/admin/dashboard");
      } else {
        navigate("/products");
      }
    } catch (err: any) {
      const errorMsg = err?.response?.data?.error?.toLowerCase() || "";
      if (errorMsg.includes("invalid credentials")) {
        setError(messages.invalidCredentials || "Invalid email or password.");
      } else {
        setError(err?.response?.data?.error || messages.genericError);
      }
    } finally {
      setLoading(false);
    }
  }

  return { login, loading, error };
}