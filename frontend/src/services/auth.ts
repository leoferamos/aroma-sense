import api from "./api";
import type { RegisterRequest, RegisterResponse, LoginRequest, LoginResponse } from "../types/auth";

export async function registerUser(data: RegisterRequest): Promise<RegisterResponse> {
  const response = await api.post<RegisterResponse>("/users/register", data);
  return response.data;
}

export async function loginUser(data: LoginRequest): Promise<LoginResponse> {
  const response = await api.post<LoginResponse>("/users/login", data);
  return response.data;
}

export async function logoutUser(): Promise<void> {
  await api.post("/users/logout");
}
