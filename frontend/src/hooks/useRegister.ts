import { useState } from "react";
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
    } catch (err: any) {
      setError(
        err?.response?.data?.error ||
        messages.registrationFailed
      );
    } finally {
      setLoading(false);
    }
  }

  return { register, loading, error, success };
}
