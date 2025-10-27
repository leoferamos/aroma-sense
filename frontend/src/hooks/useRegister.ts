import { useState } from "react";
import { isAxiosError } from "axios";
import { registerUser } from "../services/auth";
import { messages } from "../constants/messages";
import type { RegisterRequest, RegisterResponse } from "../types/auth";

export function useRegister() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");

  async function register(data: RegisterRequest) {
    setLoading(true);
    setError("");
    setSuccess("");
    try {
      const res: RegisterResponse = await registerUser(data);
      setSuccess(res.message || messages.registrationSuccess);
    } catch (err: unknown) {
      if (isAxiosError(err)) {
        const errorMsg = err.response?.data?.error?.toLowerCase() || "";
        if (errorMsg.includes("email already registered")) {
          setError(messages.emailAlreadyRegistered);
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

  return { register, loading, error, success };
}
