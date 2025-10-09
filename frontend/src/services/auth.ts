import api from "./api";
import type { RegisterRequest, RegisterResponse } from "../types/auth";

export async function registerUser(data: RegisterRequest): Promise<RegisterResponse> {
  const response = await api.post<RegisterResponse>("/users/register", data);
  return response.data;
}
