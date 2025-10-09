import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { loginUser } from "../services/auth";
import { messages } from "../constants/messages";
import type { LoginRequest, LoginResponse } from "../types/auth";

export function useLogin() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [user, setUser] = useState<LoginResponse | null>(null);
  const navigate = useNavigate();

  async function login(data: LoginRequest) {
    setLoading(true);
    setError("");
    setUser(null);
    try {
      const res: LoginResponse = await loginUser(data);
      setUser(res);
      
      // Redirect to products page on successful login
      navigate("/products");
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

  return { login, loading, error, user };
}