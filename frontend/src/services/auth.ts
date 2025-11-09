import api from "./api";
import type { 
  RegisterRequest, 
  RegisterResponse, 
  LoginRequest, 
  LoginResponse,
  RefreshResponse 
} from "../types/auth";

export async function registerUser(data: RegisterRequest): Promise<RegisterResponse> {
  const response = await api.post<RegisterResponse>("/users/register", data);
  return response.data;
}

export async function loginUser(data: LoginRequest): Promise<LoginResponse> {
  const response = await api.post<LoginResponse>("/users/login", data);
  return response.data;
}

export async function refreshToken(): Promise<RefreshResponse> {
  const response = await api.post<RefreshResponse>("/users/refresh");
  return response.data;
}

export async function logoutUser(): Promise<void> {
  await api.post("/users/logout");
}

export async function requestPasswordReset(email: string): Promise<{ message: string }> {
  const response = await api.post<{ message: string }>("/users/reset/request", { email });
  return response.data;
}

export async function confirmPasswordReset(
  email: string,
  code: string,
  newPassword: string
): Promise<{ message: string }> {
  const response = await api.post<{ message: string }>("/users/reset/confirm", {
    email,
    code,
    new_password: newPassword,
  });
  return response.data;
}
