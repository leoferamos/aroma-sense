export interface RegisterRequest {
  email: string;
  password: string;
}

export interface RegisterResponse {
  message: string;
}

export interface ErrorResponse {
  error: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface User {
  public_id: string;
  email: string;
  role: string;
  created_at: string;
}

export interface LoginResponse {
  message: string;
  access_token: string;
  user: User;
}

export interface RefreshResponse {
  message: string;
  access_token: string;
  user: User;
}
